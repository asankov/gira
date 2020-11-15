package server

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/gira-games/client/pkg/client"
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

		games, err := s.Client.GetGames(context.Background(), &client.GetGamesRequest{Token: token, ExcludeAssigned: true})
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
			Games: games.Games,
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
			s.Session.Put(r, "error", "Game is required")
			w.Header().Add("Location", "/games/add")
			w.WriteHeader(http.StatusSeeOther)
			return
		}

		if err := s.Client.LinkGameToUser(context.Background(), &client.LinkGameToUserRequest{GameID: gameID, Token: token}); err != nil {
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

		if err := s.Client.UpdateGameProgress(context.Background(), &client.UpdateGameProgressRequest{
			GameID: gameID,
			Token:  token,
			Update: client.UpdateGameProgressChange{
				Status: client.Status(status),
			},
		}); err != nil {
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
		if err := s.Client.UpdateGameProgress(context.Background(), &client.UpdateGameProgressRequest{
			GameID: gameID,
			Token:  token,
			Update: client.UpdateGameProgressChange{
				Progress: &client.UserGameProgress{
					Current: currentProgress,
					Final:   finalProgress,
				},
			},
		}); err != nil {
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

		gamesResponse, err := s.Client.GetUserGames(context.Background(), &client.GetUserGamesRequest{Token: token})
		if err != nil {
			if errors.Is(err, client.ErrNoAuthorization) {
				w.Header().Add("Location", "/users/login")
				w.WriteHeader(http.StatusSeeOther)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		statusesResponse, err := s.Client.GetStatuses(context.Background(), &client.GetStatusesRequest{Token: token})
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
			UserGames: mapToGames(gamesResponse.UserGames),
			Statuses:  statusesResponse.Statuses,
		}

		s.render(w, r, data, listGamesPage, token)
	}
}

func mapToGames(userGames map[client.Status][]*client.UserGame) []*client.UserGame {
	res := []*client.UserGame{}
	for _, v := range userGames {
		res = append(res, v...)
	}
	return res
}

func (s *Server) handleGameCreateView() authorizedHandler {
	return func(w http.ResponseWriter, r *http.Request, token string) {

		franchises := []*client.Franchise{}
		resp, err := s.Client.GetFranchises(context.Background(), &client.GetFranchisesRequest{Token: token})
		if err != nil {
			s.Log.Warnf("Error while fetching franchises: %v", err)
		} else {
			franchises = resp.Franchises
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

		if _, err := s.Client.CreateGame(context.Background(), &client.CreateGameRequest{
			Token: token,
			Game: &client.Game{
				Name:        name,
				FranchiseID: franchiseID,
			},
		}); err != nil {
			s.Session.Put(r, "error", err.Error())

			w.Header().Add("Location", "/games/new")
			w.WriteHeader(http.StatusSeeOther)
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

		if err := s.Client.DeleteUserGame(context.Background(), &client.DeleteUserGameRequest{
			GameID: gameID,
			Token:  token,
		}); err != nil {
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
		resp, err := s.Client.GetUser(context.Background(), &client.GetUserRequest{
			Token: token,
		})
		if err != nil {
			s.Log.Errorf("Error while fetching user: %v", err)
		} else {
			data.User = &client.User{
				ID:       resp.ID,
				Username: resp.Username,
				Email:    resp.Email,
			}
		}
	}

	if err := s.Renderer.Render(w, r, data, page); err != nil {
		s.Log.Errorf("Error while calling Render: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
