package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (s *server) routes() http.Handler {
	r := mux.NewRouter()

	r.HandleFunc("/games", s.getGamesHandler()).Methods(http.MethodGet)
	r.HandleFunc("/games/{id}", s.getGameByIDHandler()).Methods(http.MethodGet)
	r.HandleFunc("/games", s.createGameHandler()).Methods(http.MethodPost)

	return s.secureHeaders(r)
}
