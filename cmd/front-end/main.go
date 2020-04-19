package main

import (
	"log"
	"net/http"
	"text/template"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		t, err := template.ParseFiles("./ui/html/home.page.tmpl", "./ui/html/base.layout.tmpl")
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Internal Server Error", 500)
			return
		}

		if err := t.Execute(w, nil); err != nil {
			log.Println(err.Error())
			http.Error(w, "Internal Server Error", 500)
		}
	}).Methods(http.MethodGet)

	fileServer := http.FileServer(http.Dir("./ui/static"))
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static", fileServer))

	log.Println("Front-end listening on port 4001")
	if err := http.ListenAndServe(":4001", router); err != nil {
		log.Fatalf("error while listening: %v", err)
	}
}
