package main

import (
	"log"
	"net/http"
	"os"
)

type server struct {
	log *log.Logger
	backEndAddr string
}

func main() {
	s := &server{
		log: log.New(os.Stdout, "", log.Ldate|log.Ltime),
		// TODO: replace this with configuration
		// TODO: replace this with full-fledged client
		backEndAddr: "http://localhost:4000",
	}

	s.log.Println("Front-end listening on port 4001")
	if err := http.ListenAndServe(":4001", s.routes()); err != nil {
		log.Fatalf("error while listening: %v", err)
	}
}
