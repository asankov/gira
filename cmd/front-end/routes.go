package main

import (
	"net/http"

	"github.com/justinas/alice"

	"github.com/gorilla/mux"
)

type game struct {
	ID   string
	Name string
}

type gamesData struct {
	Games []game
	// TODO: proper data type
	User  string
	Flash string
}

func (s *server) routes() http.Handler {
	r := mux.NewRouter()

	standartMiddleware := alice.New(s.recoverPanic, s.logRequest, s.secureHeaders)
	dynamicMiddleware := alice.New(s.session.Enable)

	r.HandleFunc("/", s.homeHandler()).Methods(http.MethodGet)
	r.Handle("/games", dynamicMiddleware.Then(s.getGamesHandler())).Methods(http.MethodGet)
	r.Handle("/games/new", dynamicMiddleware.Then(s.createGameViewHandler())).Methods(http.MethodGet)
	r.Handle("/games", dynamicMiddleware.Then(s.createGameHandler())).Methods(http.MethodPost)

	fileServer := http.FileServer(http.Dir("./ui/static"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static", fileServer))

	return standartMiddleware.Then(r)
}
