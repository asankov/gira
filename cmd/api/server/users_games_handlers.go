package server

import (
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

		games, err := s.UserGamesModel.GetUserGames(user.ID)
		if err != nil {
			s.Log.Printf("error while fetching user games: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		s.respond(w, r, games, http.StatusOK)
	}
}

func userFromRequest(r *http.Request) (*models.User, error) {
	args := mux.Vars(r)
	id := args["id"]
	if id == "" {
		return nil, fmt.Errorf(`No 'id' in request`)
	}
	usr := r.Context().Value(userKey)
	if usr == nil {
		return nil, fmt.Errorf("No user found in request")
	}
	user, ok := usr.(*models.User)
	if !ok {
		return nil, fmt.Errorf("No user found in request")
	}

	if user.ID != id {
		return nil, fmt.Errorf("ID param differs from userID in token")
	}

	return user, nil
}
