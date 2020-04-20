package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

type game struct {
	ID   string
	Name string
}

type gamesData struct {
	Games []game
	// TODO: proper data type
	User string
}

func (s *server) routes() http.Handler {
	r := mux.NewRouter()

	r.HandleFunc("/", s.homeHandler()).Methods(http.MethodGet)
	r.HandleFunc("/games", s.getGamesHandler()).Methods(http.MethodGet)
	r.HandleFunc("/games/new", s.createGameViewHandler()).Methods(http.MethodGet)
	r.HandleFunc("/games", s.createGameHandler()).Methods(http.MethodPost)

	fileServer := http.FileServer(http.Dir("./ui/static"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static", fileServer))

	return r
}
