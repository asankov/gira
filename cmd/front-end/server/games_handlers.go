package server

import (
	"errors"
	"net/http"

	"github.com/asankov/gira/pkg/client"
	"github.com/asankov/gira/pkg/models"
)

func (s *Server) handleHome() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.render(w, r, emptyTemplateData, homePage)
	}
}
func (s *Server) handleGamesGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		flash := s.Session.PopString(r, "flash")

		token := getToken(r)
		gamesResponse, err := s.Client.GetUserGames(token)
		if err != nil {
			if errors.Is(err, client.ErrNoAuthorization) {
				w.Header().Add("Location", "/users/login")
				w.WriteHeader(http.StatusSeeOther)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		data := &TemplateData{
			Flash:     flash,
			UserGames: mapToGames(gamesResponse),
		}

		s.render(w, r, data, listGamesPage)
	}
}

func mapToGames(userGames map[models.Status][]*models.UserGame) []*models.UserGame {
	res := []*models.UserGame{}
	for _, v := range userGames {
		res = append(res, v...)
	}
	return res
}

func (s *Server) handleGameCreateView() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.render(w, r, emptyTemplateData, createGamePage)
	}
}

func (s *Server) handleGameCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		name := r.PostForm.Get("name")
		if name == "" {
			http.Error(w, "'name' is required", http.StatusBadRequest)
			return
		}

		token := r.Context().Value(contextTokenKey).(string)

		if _, err := s.Client.CreateGame(&models.Game{Name: name}, token); err != nil {
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

func (s *Server) render(w http.ResponseWriter, r *http.Request, data *TemplateData, p string) {
	if cookie, err := r.Cookie("token"); err == nil {
		usr, err := s.Client.GetUser(cookie.Value)
		if err != nil {
			s.Log.Printf("error while fetching user: %v", err)
		}

		data.User = usr
	}
	if err := s.Renderer.Render(w, r, data, p); err != nil {
		s.Log.Printf("error while calling Render: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
