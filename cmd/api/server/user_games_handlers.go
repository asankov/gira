package server

import (
	"encoding/json"
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
			s.Log.Printf("error while fetching user games: %v", err)
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

		var req models.UserGameRequest
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

type userGamePatchRequest struct {
	Status models.Status `json:"status"`
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

		var req userGamePatchRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			s.Log.Printf("Error while decoding body: %v", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusBadRequest)
			return
		}

		if err := s.UserGamesModel.ChangeGameStatus(user.ID, userGameID, req.Status); err != nil {
			s.Log.Printf("Error while changing game status: %v", err)
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
