// A lot of this was adapted from https://gist.github.com/ismasan/3fb75381cd2deb6bfa9c
// Great big thanks to Ismael Celis for that work.
package server

import (
	"embed"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/demmydemon/httptail/config"
)

//go:embed static/*
var static embed.FS

type TailServer struct {
	addClient    chan chan TailMessage
	removeClient chan chan TailMessage
	clients      map[chan TailMessage]time.Time
	Messaging    chan TailMessage
	useEmbedded  bool
	buffer       MessageBuffer
}

func NewServer(cfg config.Configuration) (server *TailServer) {
	server = &TailServer{
		addClient:    make(chan chan TailMessage),
		removeClient: make(chan chan TailMessage),
		clients:      make(map[chan TailMessage]time.Time),
		Messaging:    make(chan TailMessage, 1),
		useEmbedded:  cfg.UseEmbedded,
		buffer:       NewMessageBuffer(cfg.BufferLength),
	}

	go server.dispatch()
	go server.listen(cfg.Port)
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
			server.buffer.Add(message)
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

	log.Printf("--> %s", req.RemoteAddr)

	setHeaders(rw)

	clientChannel := make(chan TailMessage)
	server.addClient <- clientChannel

	defer func() {
		server.removeClient <- clientChannel
		log.Printf("<--x %s", req.RemoteAddr)
	}()

	done := req.Context().Done()
	go func() {
		<-done // That is, pause until done channel is closed
		server.removeClient <- clientChannel
		log.Printf("<-- %s", req.RemoteAddr)
	}()

	messages := server.buffer.Get()

	if len(messages) > 0 {

		connect, err := json.Marshal(TailMessage{Context: "connect", Line: fmt.Sprintf("Backbuffer is %d lines", len(server.buffer.content))})
		if err != nil {
			log.Fatalf(err.Error())
		}
		fmt.Fprintf(rw, "data: %s\n\n", connect)
		flusher.Flush()

		for _, message := range messages {
			msg, err := json.Marshal(message)
			if err != nil {
				log.Fatalf(err.Error())
			}
			fmt.Fprintf(rw, "data: %s\n\n", msg)
			flusher.Flush()
		}
	}

	for {
		message := <-clientChannel
		msg, err := json.Marshal(message)
		if err != nil {
			log.Fatalf(err.Error())
		}
		fmt.Fprintf(rw, "data: %s\n\n", msg)
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
