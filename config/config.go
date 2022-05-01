package config

import (
	"flag"
	"log"
)

type Configuration struct {
	Port         int
	Filter       string
	BufferLength int
	UseEmbedded  bool
	Files        []string
}

func GetConfiguration() Configuration {

	log.SetFlags(log.Lshortfile | log.Lmicroseconds)

	cmdFlags := Configuration{}

	flag.IntVar(&cmdFlags.Port, "port", 1030, "Port to bind to for HTTP request listening")
	flag.IntVar(&cmdFlags.BufferLength, "buffer-length", 10, "How many lines of scrollback to send new clients")
	flag.StringVar(&cmdFlags.Filter, "filter", "", "Only show lines containing this string")
	flag.BoolVar(&cmdFlags.UseEmbedded, "embedded", true, "Use the embedded HTML/CSS/JS or not")
	flag.Parse()

	cmdFlags.Files = flag.Args()

	return cmdFlags
}
