package server

import (
	"net/http"

	"github.com/asankov/gira/internal/middleware"

	"github.com/justinas/alice"

	"github.com/gorilla/mux"
)

func (s *Server) routes() http.Handler {
	standartMiddleware := alice.New(middleware.RecoverPanic(s.Log), middleware.LogRequest(s.Log))

	r := mux.NewRouter()

	// GET /games returns all games for the authorized user
	r.Handle("/games", s.requireLogin(s.handleGamesGet())).Methods(http.MethodGet)
	// POST /games creates a game for the authenticated user
	r.Handle("/games", s.requireLogin(s.handleGamesCreate())).Methods(http.MethodPost)
	// GET /games/{id} returns the requested game for the authorized user
	r.Handle("/games/{id}", s.requireLogin(s.handleGamesGetByID())).Methods(http.MethodGet)
	// PATCH /games/{id} changes the status or progress of the given games for the authenticated user
	r.Handle("/games/{id}", s.requireLogin(s.handleUsersGamesPatch())).Methods(http.MethodPatch)
	// DELETE /games/{id} deletes the given game for the authenticated user
	r.Handle("/games/{id}", s.requireLogin(s.handleUsersGamesDelete())).Methods(http.MethodDelete)

	r.HandleFunc("/users", s.handleUserGet()).Methods(http.MethodGet)
	r.HandleFunc("/users", s.handleUserCreate()).Methods(http.MethodPost)
	r.HandleFunc("/users/login", s.handleUserLogin()).Methods(http.MethodPost)

	r.Handle("/users/logout", s.requireLogin(s.handleUserLogout())).Methods(http.MethodPost)

	r.Handle("/franchises", s.requireLogin(s.handleFranchisesGet())).Methods(http.MethodGet)
	r.Handle("/franchises", s.requireLogin(s.handleFranchisesCreate())).Methods(http.MethodPost)

	r.Handle("/statuses", s.requireLogin(s.handleStatusesGet())).Methods(http.MethodGet)

	return standartMiddleware.Then(r)
}
