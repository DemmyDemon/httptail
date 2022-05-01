// A lot of this was adapted from https://gist.github.com/ismasan/3fb75381cd2deb6bfa9c
// Great big thanks to Ismael Celis for that work.
package server

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

type TailMessage struct {
	Payload []byte
}

type TailServer struct {
	addClient    chan chan TailMessage
	removeClient chan chan TailMessage
	clients      map[chan TailMessage]time.Time
	Messaging    chan TailMessage
}

func NewServer(port int) (server *TailServer) {
	server = &TailServer{
		addClient:    make(chan chan TailMessage),
		removeClient: make(chan chan TailMessage),
		clients:      make(map[chan TailMessage]time.Time),
		Messaging:    make(chan TailMessage, 1),
	}
	go server.dispatch()
	go server.listen(port)
	return
}

func (server *TailServer) listen(port int) {
	log.Fatal("Fatal error in HTTP server: ", http.ListenAndServe(fmt.Sprintf(":%d", port), server))
}

func (server *TailServer) dispatch() {
	for {
		select {
		case client := <-server.addClient:
			server.clients[client] = time.Now()
			log.Printf("Addded client, %d now connected", len(server.clients))
		case client := <-server.removeClient:
			delete(server.clients, client)
			log.Printf("Removed client, %d now connected", len(server.clients))
		case message := <-server.Messaging:
			for clientChannel := range server.clients {
				// log.Print(message.Payload)
				clientChannel <- message
			}
		}
	}
}

func (server *TailServer) ServeEvents(rw http.ResponseWriter, req *http.Request) {

	flusher, ok := rw.(http.Flusher)
	if !ok {
		http.Error(rw, "Sorry, I can't stream events to you.", http.StatusInternalServerError)
	}

	log.Println("Setting headers")
	setHeaders(rw)

	log.Println("Making client channel")
	clientChannel := make(chan TailMessage)
	server.addClient <- clientChannel

	fmt.Fprintln(rw, "Welcome aboard")
	flusher.Flush()

	defer func() {
		log.Println("Closing client channel")
		server.removeClient <- clientChannel
	}()

	done := req.Context().Done()
	go func() {
		<-done // That is, pause until done channel is closed
		log.Println("Context reports done: Closing client channel")
		server.removeClient <- clientChannel
	}()

	for {
		message := <-clientChannel
		fmt.Fprintf(rw, "data: %s\n\n", message.Payload)
		flusher.Flush()
	}

}

func setHeaders(rw http.ResponseWriter) {
	rw.Header().Set("Content-Type", "text/event-stream")
	rw.Header().Set("Connection", "keep-alive")
	rw.Header().Set("Cache-Control", "no-cache")
	rw.Header().Set("Access-Control-Allow-Origin", "*")
}

func (server *TailServer) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	log.Println(req.URL.Path)
	switch req.URL.Path {
	case "/":
		// server.serveStatic(rw, req, "static/index.html")
		log.Println("Serving HTML")
		http.ServeFile(rw, req, "server/static/index.html")
	case "/style.css":
		log.Println("Serving stylesheet")
		http.ServeFile(rw, req, "server/static/style.css")
	case "/httptail.js":
		log.Println("Serving JavaScript")
		http.ServeFile(rw, req, "server/static/httptail.js")
	case "/events":
		log.Println("Adding event client")
		server.ServeEvents(rw, req)
	default:
		http.Error(rw, "Never heard of her.", http.StatusNotFound)
	}
}
