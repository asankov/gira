package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (s *server) routes() http.Handler {
	r := mux.NewRouter()

	r.HandleFunc("/games", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("{}"))
	}).Methods(http.MethodGet)
	r.HandleFunc("/games", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("{}"))
	}).Methods(http.MethodPost)

	return r
}
