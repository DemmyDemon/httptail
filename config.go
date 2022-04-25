package main

import (
	"flag"
	"log"
	"strings"
)

type Configuration struct {
	verbose   bool
	Port      int
	FileMatch string
	Directory string
	Filter    string
}

func (cfg *Configuration) Verbose(message ...string) {
	if cfg.verbose {
		log.Println(strings.Join(message, " "))
	}
}

func GetConfiguration() Configuration {

	log.SetFlags(log.Lshortfile | log.Lmicroseconds)

	cmdFlags := Configuration{}

	flag.BoolVar(&cmdFlags.verbose, "verbose", false, "Turn on verbosity")
	flag.IntVar(&cmdFlags.Port, "port", 1030, "Port to bind to for HTTP request listening")
	flag.StringVar(&cmdFlags.FileMatch, "file-match", "", "Only tail files with names containing this")
	flag.StringVar(&cmdFlags.Directory, "directory", "", "Use this directory in place of the CWD")
	flag.StringVar(&cmdFlags.Filter, "filter", "", "Only show lines containing this string")
	flag.Parse()

	return cmdFlags
}
