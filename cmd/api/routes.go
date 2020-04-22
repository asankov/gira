package main

import (
	"net/http"

	"github.com/justinas/alice"

	"github.com/gorilla/mux"
)

func (s *server) routes() http.Handler {
	standartMiddleware := alice.New(s.recoverPanic, s.logRequest, s.secureHeaders)

	r := mux.NewRouter()

	r.HandleFunc("/games", s.getGamesHandler()).Methods(http.MethodGet)
	r.HandleFunc("/games/{id}", s.getGameByIDHandler()).Methods(http.MethodGet)
	r.HandleFunc("/games", s.createGameHandler()).Methods(http.MethodPost)

	r.HandleFunc("/users", s.createUserHandler()).Methods(http.MethodPost)
	r.HandleFunc("/users/login", s.loginHandler()).Methods(http.MethodPost)

	return standartMiddleware.Then(r)
}
