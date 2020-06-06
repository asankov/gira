package server

import (
	"net/http"

	"github.com/asankov/gira/internal/middleware"

	"github.com/justinas/alice"

	"github.com/gorilla/mux"
)

func (s *Server) routes() http.Handler {
	standartMiddleware := alice.New(middleware.RecoverPanic(s.Log), middleware.LogRequest(s.Log))
	requireLogin := alice.New(s.requireLogin)

	r := mux.NewRouter()

	r.Handle("/games", requireLogin.Then(s.handleGamesGet())).Methods(http.MethodGet)
	r.Handle("/games/{id}", requireLogin.Then(s.handleGamesGetByID())).Methods(http.MethodGet)
	r.Handle("/games", requireLogin.Then(s.handleGamesCreate())).Methods(http.MethodPost)

	r.HandleFunc("/users", s.handleUserGet()).Methods(http.MethodGet)
	r.HandleFunc("/users", s.handleUserCreate()).Methods(http.MethodPost)
	r.HandleFunc("/users/login", s.handleUserLogin()).Methods(http.MethodPost)

	r.Handle("/users/games", requireLogin.Then(s.handleUsersGamesGet())).Methods(http.MethodGet)
	r.Handle("/users/games", requireLogin.Then(s.handleUsersGamesPost())).Methods(http.MethodPost)
	r.Handle("/users/games/{id}", requireLogin.Then(s.handleUsersGamesPatch())).Methods(http.MethodPatch)

	return standartMiddleware.Then(r)
}
