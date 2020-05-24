package server

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

func (s *Server) handleHome() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.render(w, r, &gamesData{}, homePage)
	}
}
func (s *Server) handleGamesGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		flash := s.Session.PopString(r, "flash")

		token := getToken(r)
		games, err := s.Client.GetGames(token)
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

		s.render(w, r, data, listGamesPage)
	}
}

func (s *Server) handleGameCreateView() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.render(w, r, &gamesData{}, createGamePage)
	}
}

func (s *Server) handleGameCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := r.PostFormValue("name")
		if name == "" {
			http.Error(w, "'name' is required", http.StatusBadRequest)
			return
		}

		if _, err := s.Client.CreateGame(&models.Game{Name: name}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		s.Session.Put(r, "flash", "Game successfully created.")

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

func (s *Server) render(w http.ResponseWriter, r *http.Request, data interface{}, p string) {
	if err := s.Renderer.Render(w, r, data, p); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
