package server

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

type TemplateRenderer struct{}

func (t *TemplateRenderer) Render(w http.ResponseWriter, r *http.Request, d interface{}, p Page) error {
	tt, err := template.ParseFiles("./ui/html/"+string(p), "./ui/html/base.layout.tmpl")
	if err != nil {
		return err
	}

	if err := tt.Execute(w, d); err != nil {
		return err
	}
	return nil
}

func (s *Server) renderTemplate(w http.ResponseWriter, r *http.Request, data Data, templates ...string) {
	s.setUserData(data, r)
	t, err := template.ParseFiles(templates...)
	if err != nil {
		s.Log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}

	if err := t.Execute(w, data); err != nil {
		s.Log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}
}

func (s *Server) setUserData(data Data, r *http.Request) {
	cookie, err := r.Cookie("token")
	if err != nil {
		return
	}
	usr, err := s.Client.GetUser(cookie.Value)
	if err != nil {
		return
	}
	if data == nil {
		return
	}

	data.SetUser(usr)
}
