package server

import (
	"encoding/json"
	"net/http"

	"github.com/asankov/gira/pkg/models"
	"github.com/gorilla/mux"
)

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
			if err := s.GameModel.ChangeGameStatus(user.ID, userGameID, req.Status); err != nil {
				s.Log.Errorf("Error while changing game status: %v", err)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
		}

		if req.Progress != nil {
			if err := s.GameModel.ChangeGameProgress(user.ID, userGameID, req.Progress); err != nil {
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
		gameID := args["id"]
		if gameID == "" {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		if err := s.GameModel.DeleteGame(user.ID, gameID); err != nil {
			s.Log.Errorf("Error while deleting user game: %v", err)
			// TODO: determine if user does not own the game or other error
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		s.respond(w, r, nil, http.StatusOK)
	}
}
