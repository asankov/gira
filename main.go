package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type server struct {
	log     log.Logger
	handler http.Handler
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/games", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("{}"))
	})

	log.Println("listening on port 4000")
	if err := http.ListenAndServe(":4000", r); err != nil {
		log.Fatalf("error while serving: %v", err)
	}
}
