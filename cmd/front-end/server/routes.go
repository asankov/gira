package server

import (
	"net/http"

	"github.com/asankov/gira/internal/middleware"

	"github.com/justinas/alice"

	"github.com/gorilla/mux"
)

func (s *Server) routes() http.Handler {
	r := mux.NewRouter()

	// / handles the home page
	r.HandleFunc("/", s.handleHome()).Methods(http.MethodGet)

	// GET /games renders the All Games view for the authenticated user
	r.Handle("/games", s.requireLogin(s.handleGamesGetView())).Methods(http.MethodGet)

	// GET /games/new renders the New Game view for the authenticated user
	r.Handle("/games/new", s.requireLogin(s.handleGameCreateView())).Methods(http.MethodGet)
	// POST /games/new handles the creation of a new game
	r.Handle("/games/new", s.requireLogin(s.handleGameCreate())).Methods(http.MethodPost)

	r.Handle("/games/status", s.requireLogin(s.handleGamesChangeStatus())).Methods(http.MethodPost)
	r.Handle("/games/progress", s.requireLogin(s.handleGamesChangeProgress())).Methods(http.MethodPost)
	r.Handle("/games/delete", s.requireLogin(s.handleGamesDelete())).Methods(http.MethodPost)

	r.Handle("/franchises/add", s.requireLogin(s.handleFranchisesAddPost())).Methods(http.MethodPost)

	r.Handle("/users/signup", s.handleUserSignupForm()).Methods(http.MethodGet)
	r.Handle("/users/create", s.handleUserSignup()).Methods(http.MethodPost)

	r.Handle("/users/login", s.handleUserLoginForm()).Methods(http.MethodGet)
	r.Handle("/users/login", s.handleUserLogin()).Methods(http.MethodPost)
	r.Handle("/users/logout", s.requireLogin(s.handleUserLogout())).Methods(http.MethodPost)

	fileServer := http.FileServer(http.Dir("./ui/static"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static", fileServer))

	standartMiddleware := alice.New(middleware.RecoverPanic(s.Log), middleware.LogRequest(s.Log), s.secureHeaders, s.Session.Enable)
	return standartMiddleware.Then(r)
}
