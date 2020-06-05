package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/asankov/gira/pkg/models"
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
			s.Log.Printf("error while fetching user games: %v", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		// gamesResponse := models.UserGameResponse{Games: games}

		s.respond(w, r, games, http.StatusOK)
	}
}

type userGameRequest struct {
	Game *models.Game `json:"game"`
}

func (s *Server) handleUsersGamesPost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := userFromRequest(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		var req userGameRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			s.Log.Printf("Error while decoding body: %v", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusBadRequest)
			return
		}

		if err := s.UserGamesModel.LinkGameToUser(user.ID, req.Game.ID); err != nil {
			s.Log.Printf("Error while linking game %s to user %s: %v", req.Game.ID, user.ID, err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		// TODO: better response
		s.respond(w, r, nil, http.StatusOK)
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
