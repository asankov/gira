package main

import (
	"log"
	"net/http"
	"os"
)

type server struct {
	log *log.Logger
}

func main() {
	s := &server{
		log: log.New(os.Stdout, "", log.Ldate|log.Ltime),
	}

	s.log.Println("Front-end listening on port 4001")
	if err := http.ListenAndServe(":4001", s.routes()); err != nil {
		log.Fatalf("error while listening: %v", err)
	}
}
