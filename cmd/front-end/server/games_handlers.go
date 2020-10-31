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
		var token string
		if cookie, err := r.Cookie("token"); err == nil {
			token = cookie.Value
		}

		s.render(w, r, emptyTemplateData, homePage, token)
	}
}

func (s *Server) handleGamesAdd() authorizedHandler {
	return func(w http.ResponseWriter, r *http.Request, token string) {

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
		}, addGamePage, token)
	}
}

func (s *Server) handleGamesAddPost() authorizedHandler {
	return func(w http.ResponseWriter, r *http.Request, token string) {

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

func (s *Server) handleGamesChangeStatus() authorizedHandler {
	return func(w http.ResponseWriter, r *http.Request, token string) {

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

func (s *Server) handleGamesChangeProgress() authorizedHandler {
	return func(w http.ResponseWriter, r *http.Request, token string) {

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

func (s *Server) handleGamesGet() authorizedHandler {
	return func(w http.ResponseWriter, r *http.Request, token string) {

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
			UserGames: mapToGames(gamesResponse),
			Statuses:  models.AllStatuses,
		}

		s.render(w, r, data, listGamesPage, token)
	}
}

func mapToGames(userGames map[models.Status][]*models.UserGame) []*models.UserGame {
	res := []*models.UserGame{}
	for _, v := range userGames {
		res = append(res, v...)
	}
	return res
}

func (s *Server) handleGameCreateView() authorizedHandler {
	return func(w http.ResponseWriter, r *http.Request, token string) {

		franchises, err := s.Client.GetFranchises(token)
		if err != nil {
			s.Log.Warnf("Error while fetching franchises: %v", err)
			franchises = []*models.Franchise{}
		}

		selectedFranchiseIDquery, ok := r.URL.Query()["selectedFranchise"]
		var selectedFranchiseID string
		if ok && len(selectedFranchiseIDquery) > 0 {
			selectedFranchiseID = selectedFranchiseIDquery[0]
		}

		s.render(w, r, TemplateData{
			Franchises:          franchises,
			SelectedFranchiseID: selectedFranchiseID,
		}, createGamePage, token)
	}
}

func (s *Server) handleGameCreate() authorizedHandler {
	return func(w http.ResponseWriter, r *http.Request, token string) {
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

		if _, err := s.Client.CreateGame(&models.Game{Name: name, FranshiseID: franchiseID}, token); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		s.Session.Put(r, "flash", "Game successfully created.")

		w.Header().Add("Location", "/games/add")
		w.WriteHeader(http.StatusSeeOther)
	}
}

func (s *Server) handleGamesDelete() authorizedHandler {
	return func(w http.ResponseWriter, r *http.Request, token string) {

		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		gameID := r.PostForm.Get("game")
		if gameID == "" {
			http.Error(w, "'game' is required", http.StatusBadRequest)
			return
		}

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
func (s *Server) render(w http.ResponseWriter, r *http.Request, data TemplateData, page string, token string) {
	flash := s.Session.PopString(r, "flash")
	if flash != "" {
		data.Flash = flash
	}

	err := s.Session.PopString(r, "error")
	if err != "" {
		data.Error = err
	}

	if token != "" {
		usr, err := s.Client.GetUser(token)
		if err != nil {
			s.Log.Errorf("Error while fetching user: %v", err)
		} else {
			data.User = usr
		}
	}

	if err := s.Renderer.Render(w, r, data, page); err != nil {
		s.Log.Errorf("Error while calling Render: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
