package server

import (
	"net/http"

	"github.com/gira-games/api/internal/middleware"

	"github.com/justinas/alice"

	"github.com/gorilla/mux"
)

func (s *Server) routes() http.Handler {
	standartMiddleware := alice.New(middleware.RecoverPanic(s.Log), middleware.LogRequest(s.Log))

	r := mux.NewRouter()

	r.Handle("/games", s.requireLogin(s.handleGamesGet())).Methods(http.MethodGet)
	r.Handle("/games/{id}", s.requireLogin(s.handleGamesGetByID())).Methods(http.MethodGet)
	r.Handle("/games", s.requireLogin(s.handleGamesCreate())).Methods(http.MethodPost)

	r.HandleFunc("/users", s.handleUserGet()).Methods(http.MethodGet)
	r.HandleFunc("/users", s.handleUserCreate()).Methods(http.MethodPost)
	r.HandleFunc("/users/login", s.handleUserLogin()).Methods(http.MethodPost)

	r.Handle("/users/logout", s.requireLogin(s.handleUserLogout())).Methods(http.MethodPost)

	r.Handle("/users/games", s.requireLogin(s.handleUsersGamesGet())).Methods(http.MethodGet)
	r.Handle("/users/games", s.requireLogin(s.handleUsersGamesPost())).Methods(http.MethodPost)
	r.Handle("/users/games/{id}", s.requireLogin(s.handleUsersGamesPatch())).Methods(http.MethodPatch)
	r.Handle("/users/games/{id}", s.requireLogin(s.handleUsersGamesDelete())).Methods(http.MethodDelete)

	r.Handle("/franchises", s.requireLogin(s.handleFranchisesGet())).Methods(http.MethodGet)
	r.Handle("/franchises", s.requireLogin(s.handleFranchisesCreate())).Methods(http.MethodPost)

	return standartMiddleware.Then(r)
}
