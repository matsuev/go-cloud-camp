package main

import (
	"flag"
	"fmt"
	"go-cloud-camp/server"
	"log"
)

// Default path to config file
const (
	DEFAULT_CONFIG_PATH = "./config.yml"
)

func main() {
	fmt.Println("GoCloudCamp config server")

	// Set path to config file from command line.
	cfgPath := flag.String("c", DEFAULT_CONFIG_PATH, "path to config file")
	flag.Parse()

	srv, err := server.Create(*cfgPath)
	if err != nil {
		log.Fatalln("Couldn't start server, caused error:", err)
	}

	srv.Run()
}
