package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"text/template"

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

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		s.renderTemplate(w, r, nil, "./ui/html/home.page.tmpl", "./ui/html/base.layout.tmpl")
	}).Methods(http.MethodGet)

	// TODO: fetch this from the back-end, instead of hardcoding them
	data := &gamesData{
		Games: []game{
			{ID: "1", Name: "AC"},
			{ID: "2", Name: "ACII"},
			{ID: "3", Name: "ACII: Brotherhood"},
			{ID: "4", Name: "ACII: Liberation"},
		},
	}
	r.HandleFunc("/games", func(w http.ResponseWriter, r *http.Request) {
		s.renderTemplate(w, r, data, "./ui/html/list.page.tmpl", "./ui/html/base.layout.tmpl")
	}).Methods(http.MethodGet)

	r.HandleFunc("/games/new", func(w http.ResponseWriter, r *http.Request) {
		s.renderTemplate(w, r, nil, "./ui/html/create.page.tmpl", "./ui/html/base.layout.tmpl")
	}).Methods(http.MethodGet)

	r.HandleFunc("/games", func(w http.ResponseWriter, r *http.Request) {
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

		w.WriteHeader(http.StatusOK)
		w.Write(b)
	}).Methods(http.MethodPost)

	fileServer := http.FileServer(http.Dir("./ui/static"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static", fileServer))

	return r
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
