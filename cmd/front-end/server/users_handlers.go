package server

import (
	"net/http"

	"github.com/asankov/gira/pkg/client"
	"github.com/asankov/gira/pkg/models"
)

func (s *Server) handleUserSignupForm() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.render(w, r, emptyTemplateData, signupUserPage, "")
	}
}

func (s *Server) handleUserLoginForm() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.render(w, r, emptyTemplateData, loginUserPage, "")
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
			s.Log.Errorf("Error while logging in user: %v", err)
			// TODO: return to login screen with sensible error
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "token",
			Value:    res.Token,
			Path:     "/",
			SameSite: http.SameSiteLaxMode,
		})
		w.Header().Add("Location", "/")
		w.WriteHeader(http.StatusSeeOther)
	}
}

func (s *Server) handleUserLogout() authorizedHandler {
	return func(w http.ResponseWriter, r *http.Request, token string) {
		if err := s.Client.LogoutUser(token); err != nil {
			// TODO: render error page
			s.Log.Printf("Error while logging-out user: %v", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

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

		email, password := r.PostForm.Get("email"), r.PostForm.Get("password")
		if email == "" || password == "" {
			http.Error(w, "email and password are required", http.StatusBadRequest)
			return
		}

		if _, err := s.Client.CreateUser(&models.User{
			Email:    email,
			Password: password,
		}); err != nil {
			s.Log.Errorf("Error while creating user: %v %v", err, err == nil)
			if errResponse, ok := err.(*client.ErrorResponse); ok {
				s.render(w, r, TemplateData{Error: errResponse.Error()}, signupUserPage, "")
				return
			}
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Add("Location", "/")
		w.WriteHeader(http.StatusSeeOther)
	}
}
