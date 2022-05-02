package main

import (
	"os"
	"os/signal"

	"github.com/demmydemon/httptail/config"
	"github.com/demmydemon/httptail/server"
	"github.com/demmydemon/httptail/tailer"
)

func main() {
	cfg := config.GetConfiguration()
	srv := server.NewServer(cfg)
	tailer.TailFiles(cfg, srv)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}
