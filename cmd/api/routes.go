package main

import (
	"github.com/justinas/alice"
	"net/http"

	"github.com/gorilla/mux"
)

func (s *server) routes() http.Handler {
	standartMiddleware := alice.New(s.recoverPanic, s.logRequest, s.secureHeaders)

	r := mux.NewRouter()

	r.HandleFunc("/games", s.getGamesHandler()).Methods(http.MethodGet)
	r.HandleFunc("/games/{id}", s.getGameByIDHandler()).Methods(http.MethodGet)
	r.HandleFunc("/games", s.createGameHandler()).Methods(http.MethodPost)

	return standartMiddleware.Then(r)
}
