package server

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/asankov/gira/pkg/client"
	"github.com/asankov/gira/pkg/models"
)

func (s *Server) handleHome() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.render(w, r, emptyTemplateData, homePage)
	}
}

func (s *Server) handleGamesAdd() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		token := getToken(r)
		games, err := s.Client.GetGames(token, &client.GetGamesOptions{ExcludeAssigned: true})
		if err != nil {
			if errors.Is(err, client.ErrNoAuthorization) {
				w.Header().Add("Location", "/users/login")
				w.WriteHeader(http.StatusSeeOther)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		s.render(w, r, TemplateData{
			Games: games,
		}, addGamePage)
	}
}

func (s *Server) handleGamesAddPost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := getToken(r)

		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		gameID := r.PostForm.Get("game")
		if gameID == "" {
			http.Error(w, "'game' is required", http.StatusBadRequest)
			return
		}

		if _, err := s.Client.LinkGameToUser(gameID, token); err != nil {
			// TODO: if err == no auth
			s.Log.Errorln(err)
			// TODO: render error page
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Add("Location", "/games")
		w.WriteHeader(http.StatusSeeOther)
	}
}

func (s *Server) handleGamesChangeStatus() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := getToken(r)

		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		gameID := r.PostForm.Get("game")
		if gameID == "" {
			http.Error(w, "'game' is required", http.StatusBadRequest)
			return
		}

		status := r.PostForm.Get("status")
		if status == "" {
			http.Error(w, "'status' is requred", http.StatusBadRequest)
			return
		}

		if err := s.Client.ChangeGameStatus(gameID, token, models.Status(status)); err != nil {
			s.Log.Errorln(err)
			// TODO: render error page
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Add("Location", "/games")
		w.WriteHeader(http.StatusSeeOther)
	}
}

func (s *Server) handleGamesChangeProgress() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := getToken(r)

		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		gameID := r.PostForm.Get("game")
		if gameID == "" {
			http.Error(w, "'game' is required", http.StatusBadRequest)
			return
		}

		cur, fin := r.PostForm.Get("currentProgress"), r.PostForm.Get("finalProgress")
		if cur == "" || fin == "" {
			http.Error(w, "'currentProgress' and 'finalProgress' is required", http.StatusBadRequest)
			return
		}

		currentProgress, err := strconv.Atoi(cur)
		if err != nil {
			http.Error(w, "'currentProgress' should be a valid integer", http.StatusBadRequest)
			return
		}
		finalProgress, err := strconv.Atoi(fin)
		if err != nil {
			http.Error(w, "'finalProgress' should be a valid integer", http.StatusBadRequest)
			return
		}
		if err := s.Client.ChangeGameProgress(gameID, token, &models.UserGameProgress{Current: currentProgress, Final: finalProgress}); err != nil {
			s.Log.Errorln(err)
			// TODO: render error page
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Add("Location", "/games")
		w.WriteHeader(http.StatusSeeOther)
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

		data := TemplateData{
			Flash:     flash,
			UserGames: mapToGames(gamesResponse),
			Statuses:  models.AllStatuses,
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
		token := getToken(r)

		franchises, err := s.Client.GetFranchises(token)
		if err != nil {
			s.Log.Warnf("Error while fetching franchises: %v", err)
			franchises = []*models.Franchise{}
		}

		selectedFranchiseIDquery, _ := r.URL.Query()["selectedFranchise"]
		var selectedFranchiseID string
		if len(selectedFranchiseIDquery) > 0 {
			selectedFranchiseID = selectedFranchiseIDquery[0]
		}

		s.render(w, r, TemplateData{
			Franchises:          franchises,
			SelectedFranchiseID: selectedFranchiseID,
		}, createGamePage)
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

		franchiseID := r.PostForm.Get("franchiseId")

		token := getToken(r)

		if _, err := s.Client.CreateGame(&models.Game{Name: name, FranshiseID: franchiseID}, token); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		s.Session.Put(r, "flash", "Game successfully created.")

		w.Header().Add("Location", "/games/add")
		w.WriteHeader(http.StatusSeeOther)
	}
}

func (s *Server) handleGamesDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		gameID := r.PostForm.Get("game")
		if gameID == "" {
			http.Error(w, "'game' is required", http.StatusBadRequest)
			return
		}

		token := getToken(r)

		if err := s.Client.DeleteUserGame(gameID, token); err != nil {
			if errors.Is(err, client.ErrNoAuthorization) {
				w.Header().Add("Location", "/users/login")
				w.WriteHeader(http.StatusSeeOther)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		s.Session.Put(r, "flash", "Game successfully deleted.")

		w.Header().Add("Location", "/games")
		w.WriteHeader(http.StatusSeeOther)
	}
}

func getToken(r *http.Request) string {
	val := r.Context().Value(contextTokenKey)
	if token, ok := val.(string); ok {
		return token
	}
	return ""
}

func (s *Server) render(w http.ResponseWriter, r *http.Request, data TemplateData, p string) {
	if token := getToken(r); token != "" {
		usr, err := s.Client.GetUser(token)
		if err != nil {
			s.Log.Errorf("Error while fetching user: %v", err)
		} else {
			data.User = usr
		}
	}

	if err := s.Renderer.Render(w, r, data, p); err != nil {
		s.Log.Errorf("Error while calling Render: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
