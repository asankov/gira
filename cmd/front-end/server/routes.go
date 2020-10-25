package server

import (
	"net/http"

	"github.com/asankov/gira/internal/middleware"

	"github.com/justinas/alice"

	"github.com/gorilla/mux"
)

func (s *Server) routes() http.Handler {
	r := mux.NewRouter()

	dynamicMiddleware := alice.New(s.Session.Enable)

	r.HandleFunc("/", s.handleHome()).Methods(http.MethodGet)
	r.Handle("/games", dynamicMiddleware.Then(s.requireLogin(s.handleGamesGet()))).Methods(http.MethodGet)
	r.Handle("/games/add", dynamicMiddleware.Then(s.requireLogin(s.handleGamesAdd()))).Methods(http.MethodGet)
	r.Handle("/games/add", dynamicMiddleware.Then(s.requireLogin(s.handleGamesAddPost()))).Methods(http.MethodPost)
	r.Handle("/games/new", dynamicMiddleware.Then(s.requireLogin(s.handleGameCreateView()))).Methods(http.MethodGet)
	r.Handle("/games", dynamicMiddleware.Then(s.requireLogin(s.handleGameCreate()))).Methods(http.MethodPost)
	r.Handle("/games/status", dynamicMiddleware.Then(s.requireLogin(s.handleGamesChangeStatus()))).Methods(http.MethodPost)
	r.Handle("/games/progress", dynamicMiddleware.Then(s.requireLogin(s.handleGamesChangeProgress()))).Methods(http.MethodPost)
	r.Handle("/games/delete", dynamicMiddleware.Then(s.requireLogin(s.handleGamesDelete()))).Methods(http.MethodPost)

	r.Handle("/franchises/add", dynamicMiddleware.Then(s.requireLogin(s.handleFranchisesAddPost()))).Methods(http.MethodPost)

	r.Handle("/users/signup", dynamicMiddleware.Then(s.handleUserSignupForm())).Methods(http.MethodGet)
	r.Handle("/users/create", dynamicMiddleware.Then(s.handleUserSignup())).Methods(http.MethodPost)

	r.Handle("/users/login", s.handleUserLoginForm()).Methods(http.MethodGet)
	r.Handle("/users/login", s.handleUserLogin()).Methods(http.MethodPost)
	r.Handle("/users/logout", s.requireLogin(s.handleUserLogout())).Methods(http.MethodPost)

	fileServer := http.FileServer(http.Dir("./ui/static"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static", fileServer))

	standartMiddleware := alice.New(middleware.RecoverPanic(s.Log), middleware.LogRequest(s.Log), s.secureHeaders)
	return standartMiddleware.Then(r)
}
