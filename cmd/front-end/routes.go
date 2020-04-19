package main

import (
	"net/http"
	"text/template"

	"github.com/gorilla/mux"
)

func (s *server) routes() http.Handler {
	r := mux.NewRouter()

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		s.renderTemplate(w, r, "./ui/html/home.page.tmpl", "./ui/html/base.layout.tmpl")
	}).Methods(http.MethodGet)

	fileServer := http.FileServer(http.Dir("./ui/static"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static", fileServer))

	return r
}

func (s *server) renderTemplate(w http.ResponseWriter, r *http.Request, templates ...string) {
	t, err := template.ParseFiles(templates...)
	if err != nil {
		s.log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}

	if err := t.Execute(w, nil); err != nil {
		s.log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}
}
