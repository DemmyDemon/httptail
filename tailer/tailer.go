package tailer

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/demmydemon/httptail/config"
	"github.com/demmydemon/httptail/server"
)

func TailFiles(cfg config.Configuration, srv *server.TailServer) {
	//TODO: Loop through  cfg.Files and attach actual tails to them.
	//TODO: Respect cfg.BufferLength
	//TODO: Send a zero-length payload every few seconds for keep-alive reasons?

	if len(cfg.Files) < 1 {
		log.Fatal("You can't tail *zero* files.")
	}

	//TODO: Support more than one file!
	if len(cfg.Files) > 1 {
		log.Println("Only the first file will be tailed, sorry! It's on the TODO list.")
	}

	file, err := os.Open(cfg.Files[0])

	if err != nil {
		log.Fatal(err.Error())
	}
	lines := FileTail(file, cfg.BufferLength)
	file.Close()
	for _, line := range lines {
		srv.Messaging <- server.TailMessage{Line: line}
	}

	go func() {
		for {
			time.Sleep(time.Second * 5)
			message := server.TailMessage{
				Line: fmt.Sprintf("Test data: %v -- Really long line so I can test the line wrap stuff! It has to be at least this long: ---------------->", time.Now()),
			}
			srv.Messaging <- message
		}
	}()

}

func FileTail(file *os.File, numberOfLines int) []string {
	stat, _ := file.Stat()
	size := stat.Size()
	var cursor int64 = 0

	char := make([]byte, 1)
	lines := []string{}
	thisLine := []byte{}

	for {

		cursor--
		file.Seek(cursor, io.SeekEnd)
		file.Read(char)

		if char[0] == 10 { // char(10) is a newline
			if cursor == -1 {
				continue // File ends with a newline
			}
			ReverseSlice(thisLine)
			lines = append(lines, string(thisLine))
			if len(lines) == numberOfLines {
				break
			}
			thisLine = []byte{}
		} else {
			thisLine = append(thisLine, char[0])
		}

		if cursor == -size { // The cursor has arrived at the beginning of the file
			ReverseSlice(thisLine)
			lines = append(lines, string(thisLine))
			break
		}

	}

	ReverseSlice(lines)
	return lines
}

// ReverseSlice rearanges a slice to be ordered in the opposite direction.
func ReverseSlice[S ~[]E, E any](s S) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}
