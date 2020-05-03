package main

import (
	"net/http"

	"github.com/justinas/alice"

	"github.com/gorilla/mux"
)

func (s *server) routes() http.Handler {
	standartMiddleware := alice.New(s.recoverPanic, s.logRequest, s.secureHeaders)

	r := mux.NewRouter()

	r.HandleFunc("/games", s.handleGamesGet()).Methods(http.MethodGet)
	r.HandleFunc("/games/{id}", s.handleGamesGetByID()).Methods(http.MethodGet)
	r.HandleFunc("/games", s.handleGamesCreate()).Methods(http.MethodPost)

	r.HandleFunc("/users", s.handleUserCreate()).Methods(http.MethodPost)
	r.HandleFunc("/users/login", s.handleUserLogin()).Methods(http.MethodPost)

	return standartMiddleware.Then(r)
}
