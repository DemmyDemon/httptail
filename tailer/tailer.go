package tailer

import (
	"fmt"
	"time"

	"github.com/demmydemon/httptail/config"
	"github.com/demmydemon/httptail/server"
)

func TailFiles(cfg config.Configuration, srv *server.TailServer) {
	//TODO: Loop through  cfg.Files and attach actual tails to them.
	//TODO: Respect cfg.BufferLength
	//TODO: Send a zero-length payload every few seconds for keep-alive reasons?

	go func() {
		for {
			time.Sleep(time.Second * 1)
			message := server.TailMessage{
				Payload: []byte(fmt.Sprintf("Test data: %v", time.Now())),
			}
			srv.Messaging <- message
		}
	}()

}
