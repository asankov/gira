package main

import (
	"errors"
	"net/http"

	"github.com/asankov/gira/pkg/client"
	"github.com/asankov/gira/pkg/models"
)

type gamesData struct {
	Games []*models.Game
	User  *models.User
	Flash string
}

// SetUser implements the Data interface
func (g *gamesData) SetUser(usr *models.User) {
	g.User = usr
}

func (s *server) handleHome() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.renderTemplate(w, r, &gamesData{}, "./ui/html/home.page.tmpl", "./ui/html/base.layout.tmpl")
	}
}
func (s *server) handleGamesGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		flash := s.session.PopString(r, "flash")

		token := getToken(r)
		games, err := s.client.GetGames(token)
		if err != nil {
			if errors.Is(err, client.ErrNoAuthorization) {
				w.Header().Add("Location", "/users/login")
				w.WriteHeader(http.StatusSeeOther)
				return
			}
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

func (s *server) handleGameCreateView() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.renderTemplate(w, r, &gamesData{}, "./ui/html/create.page.tmpl", "./ui/html/base.layout.tmpl")
	}
}

func (s *server) handleGameCreate() http.HandlerFunc {
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

func getToken(r *http.Request) string {
	cookie, err := r.Cookie("token")
	if err != nil {
		// let it panic, the middleware should not allow this to happen
		panic("token not present in cookie")
	}
	return cookie.Value
}
