package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
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

		// TODO: fetch this from the back-end, instead of hardcoding them
		res, err := http.Get(fmt.Sprintf("%s/games", s.backEndAddr))
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		if res.StatusCode != 200 {
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}

		var games []models.Game
		json.NewDecoder(res.Body).Decode(&games)

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

		resp, err := http.Post(fmt.Sprintf("%s/games", s.backEndAddr), "application/json", strings.NewReader(fmt.Sprintf(`{"name": "%s"}`, name)))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		s.session.Put(r, "flash", "Game successfully created.")

		w.WriteHeader(http.StatusOK)
		w.Write(b)
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
