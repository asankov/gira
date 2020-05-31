package server

import (
	"net/http"

	"github.com/asankov/gira/pkg/client"
	"github.com/asankov/gira/pkg/models"
)

func (s *Server) handleUserSignupForm() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.render(w, r, emptyTemplateData, signupUserPage)
	}
}

func (s *Server) handleUserLoginForm() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		nextQuery := r.URL.Query()["next"]
		next := ""
		if len(nextQuery) > 0 {
			next = "?next=" + nextQuery[0]
		}
		s.render(w, r, &TemplateData{Next: next}, loginUserPage)
	}
}

func (s *Server) handleUserLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		email, password := r.PostForm.Get("email"), r.PostForm.Get("password")
		res, err := s.Client.LoginUser(&models.User{
			Email:    email,
			Password: password,
		})
		if err != nil {
			s.Log.Printf("error while logging in user: %v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:  "token",
			Value: res.Token,
			Path:  "/",
		})

		w.Header().Add("Location", "/")
		w.WriteHeader(http.StatusSeeOther)
	}
}

func (s *Server) handleUserSignup() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		username, email, password := r.PostForm.Get("username"), r.PostForm.Get("email"), r.PostForm.Get("password")

		if _, err := s.Client.CreateUser(&models.User{
			Username: username,
			Email:    email,
			Password: password,
		}); err != nil {
			s.Log.Printf("error while creating user: %v %v", err, err == nil)
			if errResponse, ok := err.(*client.ErrorResponse); ok {
				s.render(w, r, &TemplateData{Error: errResponse.Error()}, signupUserPage)
				return
			}
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Add("Location", "/")
		w.WriteHeader(http.StatusSeeOther)
	}
}
