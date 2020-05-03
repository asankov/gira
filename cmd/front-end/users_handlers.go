package main

import (
	"net/http"

	"github.com/asankov/gira/pkg/models"
)

func (s *server) getSignupFormHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.renderTemplate(w, r, nil, "./ui/html/signup.page.tmpl", "./ui/html/base.layout.tmpl")
	}
}

func (s *server) createUserHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username, email, password := r.PostFormValue("username"), r.PostFormValue("email"), r.PostFormValue("password")

		if _, err := s.client.CreateUser(&models.User{
			Username: username,
			Email:    email,
			Password: password,
		}); err != nil {
			s.log.Printf("error while creating user: %v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// TODO: this is not shown
		s.session.Put(r, "flash", "User created succesfully.")

		w.Header().Add("Location", "/")
		w.WriteHeader(http.StatusSeeOther)
	}
}
