package main

import (
	"net/http"
	"text/template"

	"github.com/asankov/gira/pkg/models"
)

// Data is the generic interface that is accepted, by the method that renders the templates.
// Its only method is User, which accepts a reference to a models.User.
// This is merely a setter to be used by the renderTemplate method to set the user, if such is present in the request.
type Data interface {
	SetUser(*models.User)
}

func (s *server) renderTemplate(w http.ResponseWriter, r *http.Request, data Data, templates ...string) {
	cookie, err := r.Cookie("token")
	if err != nil {

	} else {
		usr, err := s.client.GetUser(cookie.Value)
		if err != nil {

		} else {
			data.SetUser(usr)
		}
	}
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
