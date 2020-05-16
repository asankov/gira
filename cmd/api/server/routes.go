package server

import (
	"net/http"

	"github.com/justinas/alice"

	"github.com/gorilla/mux"
)

func (s *Server) routes() http.Handler {
	standartMiddleware := alice.New(s.recoverPanic, s.logRequest, s.secureHeaders)
	requireLogin := alice.New(s.requireLogin)

	r := mux.NewRouter()

	r.Handle("/games", requireLogin.Then(s.handleGamesGet())).Methods(http.MethodGet)
	r.Handle("/games/{id}", requireLogin.Then(s.handleGamesGetByID())).Methods(http.MethodGet)
	r.Handle("/games", requireLogin.Then(s.handleGamesCreate())).Methods(http.MethodPost)

	r.HandleFunc("/users", s.handleUserCreate()).Methods(http.MethodPost)
	r.HandleFunc("/users/login", s.handleUserLogin()).Methods(http.MethodPost)

	return standartMiddleware.Then(r)
}
