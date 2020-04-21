package main

import (
	"net/http"
	"text/template"

	"github.com/asankov/gira/pkg/models"
)

func (s *server) homeHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.renderTemplate(w, r, nil, "./ui/html/home.page.tmpl", "./ui/html/base.layout.tmpl")
	}
}
func (s *server) getGamesHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		flash := s.session.PopString(r, "flash")

		games, err := s.client.GetGames()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		data := &gamesData{
			Flash: flash,
			Games: games,
		}

		s.renderTemplate(w, r, data, "./ui/html/list.page.tmpl", "./ui/html/base.layout.tmpl")
	}
}

func (s *server) createGameViewHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.renderTemplate(w, r, nil, "./ui/html/create.page.tmpl", "./ui/html/base.layout.tmpl")
	}
}

func (s *server) createGameHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := r.PostFormValue("name")
		if name == "" {
			http.Error(w, "'name' is required", http.StatusBadRequest)
			return
		}

		if _, err := s.client.CreateGame(&models.Game{Name: name}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		s.session.Put(r, "flash", "Game successfully created.")

		w.Header().Add("Location", "/games")
		w.WriteHeader(http.StatusSeeOther)
	}
}

func (s *server) renderTemplate(w http.ResponseWriter, r *http.Request, data interface{}, templates ...string) {
	t, err := template.ParseFiles(templates...)
	if err != nil {
		s.log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}

	if err := t.Execute(w, data); err != nil {
		s.log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}
}
