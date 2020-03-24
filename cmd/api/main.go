package main

import (
	"log"
	"net/http"
)

type server struct {
	log     log.Logger
	handler http.Handler
}

func main() {
	s := &server{}

	log.Println("listening on port 4000")
	if err := http.ListenAndServe(":4000", s.routes()); err != nil {
		log.Fatalf("error while serving: %v", err)
	}
}
