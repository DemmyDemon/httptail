// A lot of this was adapted from https://gist.github.com/ismasan/3fb75381cd2deb6bfa9c
// Great big thanks to Ismael Celis for that work.
package server

import (
	"embed"
	"fmt"
	"log"
	"net/http"
	"time"
)

//go:embed static/*
var static embed.FS

type TailMessage struct {
	Payload []byte
}

type TailServer struct {
	addClient    chan chan TailMessage
	removeClient chan chan TailMessage
	clients      map[chan TailMessage]time.Time
	Messaging    chan TailMessage
	useEmbedded  bool
}

func NewServer(port int, useEmbedded bool) (server *TailServer) {
	server = &TailServer{
		addClient:    make(chan chan TailMessage),
		removeClient: make(chan chan TailMessage),
		clients:      make(map[chan TailMessage]time.Time),
		Messaging:    make(chan TailMessage, 1),
		useEmbedded:  useEmbedded,
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
				clientChannel <- message
			}
		}
	}
}

func (server *TailServer) ServeEvents(rw http.ResponseWriter, req *http.Request) {

	flusher, ok := rw.(http.Flusher)
	if !ok {
		http.Error(rw, "Sorry, I can't stream events to you.", http.StatusInternalServerError)
		return
	}

	setHeaders(rw)

	clientChannel := make(chan TailMessage)
	server.addClient <- clientChannel

	defer func() {
		server.removeClient <- clientChannel
	}()

	done := req.Context().Done()
	go func() {
		<-done // That is, pause until done channel is closed
		server.removeClient <- clientChannel
	}()

	fmt.Fprint(rw, "data: Welcome aboard\n\n")
	flusher.Flush()

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

func (server *TailServer) ServeEmbed(rw http.ResponseWriter, req *http.Request, filename string, mimeType string) {

	if server.useEmbedded {
		rw.Header().Set("Content-Type", mimeType)
		data, err := static.ReadFile("static/" + filename)
		if err != nil {
			http.Error(rw, "Not an embedded file: "+filename, http.StatusNotFound)
		}
		rw.WriteHeader(http.StatusOK)
		rw.Write(data)
		return
	}
	http.ServeFile(rw, req, "server/static/"+filename)
}

func (server *TailServer) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	switch req.URL.Path {
	case "/":
		server.ServeEmbed(rw, req, "index.html", "text/html")
	case "/style.css":
		server.ServeEmbed(rw, req, "style.css", "text/css")
	case "/httptail.js":
		server.ServeEmbed(rw, req, "httptail.js", "text/javascript")
	case "/favicon.ico":
		server.ServeEmbed(rw, req, "favicon.ico", "image/x-icon")
	case "/events":
		server.ServeEvents(rw, req)
	default:
		http.Error(rw, "Never heard of her.", http.StatusNotFound)
	}
}
