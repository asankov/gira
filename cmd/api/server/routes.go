package server

import (
	"net/http"

	"github.com/justinas/alice"

	"github.com/gorilla/mux"
)

func (s *Server) routes() http.Handler {
	standartMiddleware := alice.New(s.recoverPanic, s.logRequest)
	requireLogin := alice.New(s.requireLogin)

	r := mux.NewRouter()

	r.Handle("/games", requireLogin.Then(s.handleGamesGet())).Methods(http.MethodGet)
	r.Handle("/games/{id}", requireLogin.Then(s.handleGamesGetByID())).Methods(http.MethodGet)
	r.Handle("/games", requireLogin.Then(s.handleGamesCreate())).Methods(http.MethodPost)

	r.HandleFunc("/users", s.handleUserGet()).Methods(http.MethodGet)
	r.HandleFunc("/users", s.handleUserCreate()).Methods(http.MethodPost)
	r.HandleFunc("/users/login", s.handleUserLogin()).Methods(http.MethodPost)

	r.HandleFunc("/users/{id}/games", nil).Methods(http.MethodGet)
	r.HandleFunc("/users/{id}/games", nil).Methods(http.MethodPost)

	return standartMiddleware.Then(r)
}
