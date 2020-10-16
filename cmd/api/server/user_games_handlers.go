package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/asankov/gira/pkg/models"
	"github.com/gorilla/mux"
)

func (s *Server) handleUsersGamesGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := userFromRequest(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

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

func (s *Server) handleUsersGamesPost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := userFromRequest(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

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

func (s *Server) handleUsersGamesPatch() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := userFromRequest(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

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

func (s *Server) handleUsersGamesDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := userFromRequest(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

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

		s.respondError(w, r, errors.New(http.StatusText(http.StatusBadRequest)), http.StatusBadRequest)
	}
}

func userFromRequest(r *http.Request) (*models.User, error) {
	usr := r.Context().Value(contextUserKey)
	if usr == nil {
		return nil, fmt.Errorf("No user found in request")
	}
	user, ok := usr.(*models.User)
	if !ok {
		return nil, fmt.Errorf("No user found in request")
	}

	return user, nil
}

func tokenFromRequest(r *http.Request) (string, error) {
	tkn := r.Context().Value(contextTokenKey)
	if tkn == nil {
		return "", fmt.Errorf("No token found in request")
	}
	token, ok := tkn.(string)
	if !ok {
		return "", fmt.Errorf("No token found in request")
	}

	return token, nil
}
