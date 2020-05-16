package main

import (
	"net/http"

	"github.com/asankov/gira/pkg/models"
	"github.com/justinas/alice"

	"github.com/gorilla/mux"
)

type gamesData struct {
	Games []*models.Game
	// TODO: proper data type
	User  string
	Flash string
}

func (s *server) routes() http.Handler {
	r := mux.NewRouter()

	standartMiddleware := alice.New(s.recoverPanic, s.logRequest, s.secureHeaders)
	dynamicMiddleware := alice.New(s.session.Enable)
	requireLogin := alice.New(s.requireLogin)

	r.HandleFunc("/", s.handleHome()).Methods(http.MethodGet)
	r.Handle("/games", requireLogin.Then(dynamicMiddleware.Then(s.handleGamesGet()))).Methods(http.MethodGet)
	r.Handle("/games/new", requireLogin.Then(dynamicMiddleware.Then(s.handleGameCreateView()))).Methods(http.MethodGet)
	r.Handle("/games", requireLogin.Then(dynamicMiddleware.Then(s.handleGameCreate()))).Methods(http.MethodPost)

	r.Handle("/users/signup", standartMiddleware.Then(s.handleUserSignupForm())).Methods(http.MethodGet)
	r.Handle("/users/create", dynamicMiddleware.Then(s.handleUserSignup())).Methods(http.MethodPost)

	r.Handle("/users/login", s.handleUserLoginForm()).Methods(http.MethodGet)
	r.Handle("/users/login", s.handleUserLogin()).Methods(http.MethodPost)

	fileServer := http.FileServer(http.Dir("./ui/static"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static", fileServer))

	return standartMiddleware.Then(r)
}
