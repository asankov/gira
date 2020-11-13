package server

import (
	"encoding/json"
	"net/http"

	"github.com/gira-games/api/pkg/models"
	"github.com/gorilla/mux"
)

func (s *Server) handleUsersGamesGet() authorizedHandler {
	return func(w http.ResponseWriter, r *http.Request, user *models.User, token string) {

		games, err := s.UserGamesModel.GetUserGamesGrouped(user.ID)
		if err != nil {
			s.Log.Errorf("Error while fetching user games: %v", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		// gamesResponse := models.UserGameResponse{Games: games}

		s.respond(w, r, games, http.StatusOK)
	}
}

func (s *Server) handleUsersGamesPost() authorizedHandler {
	return func(w http.ResponseWriter, r *http.Request, user *models.User, token string) {

		req := models.UserGameRequest{
			Progress: &models.UserGameProgress{
				Current: 0,
				Final:   100,
			},
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			s.Log.Errorf("Error while decoding body: %v", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusBadRequest)
			return
		}

		if err := s.UserGamesModel.LinkGameToUser(user.ID, req.Game.ID, req.Progress); err != nil {
			s.Log.Errorf("Error while linking game %s to user %s: %v", req.Game.ID, user.ID, err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		// TODO: better response
		s.respond(w, r, nil, http.StatusOK)
	}
}

func (s *Server) handleUsersGamesPatch() authorizedHandler {
	return func(w http.ResponseWriter, r *http.Request, user *models.User, token string) {

		args := mux.Vars(r)
		userGameID := args["id"]
		if userGameID == "" {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		var req models.ChangeGameStatusRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			s.Log.Errorf("Error while decoding body: %v", err)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		if req.Status != "" {
			if err := req.Status.Validate(); err != nil {
				s.respondError(w, r, err.Error(), http.StatusBadRequest)
				return
			}
			if err := s.UserGamesModel.ChangeGameStatus(user.ID, userGameID, req.Status); err != nil {
				s.Log.Errorf("Error while changing game status: %v", err)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
		}

		if req.Progress != nil {
			if err := s.UserGamesModel.ChangeGameProgress(user.ID, userGameID, req.Progress); err != nil {
				s.Log.Errorf("Error while changing game progress: %v", err)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
		}

		// TODO: better response
		s.respond(w, r, nil, http.StatusOK)
	}
}

func (s *Server) handleUsersGamesDelete() authorizedHandler {
	return func(w http.ResponseWriter, r *http.Request, user *models.User, token string) {

		args := mux.Vars(r)
		userGameID := args["id"]
		if userGameID == "" {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		userGames, _ := s.UserGamesModel.GetUserGames(user.ID)
		for _, userGame := range userGames {
			if userGame.ID == userGameID {
				if err := s.UserGamesModel.DeleteUserGame(userGameID); err != nil {
					s.Log.Errorf("Error while deleting user game: %v", err)
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
					return
				}

				// TODO: better response
				s.respond(w, r, nil, http.StatusOK)
			}
		}

		s.respondError(w, r, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}
}
