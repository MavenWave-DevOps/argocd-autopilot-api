package main

import (
	"github.com/tony-mw/argocd-api/server"
	"log"
)

const (
	hostname = "localhost"
	port     = 8080
)

func main() {
	s := server.Server{
		Hostname: hostname,
		Port:     port,
	}

	if err := s.Start(); err != nil {
		log.Fatal(err, "couldn't start server")
	}
}
