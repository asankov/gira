package server

import (
	"net/http"

	"github.com/asankov/gira/pkg/models"
	"github.com/gorilla/mux"
)

func (s *Server) handleUsersGamesGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		args := mux.Vars(r)
		id := args["id"]
		if id == "" {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		user := userFromRequest(r)

		if user.ID != id {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
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

func userFromRequest(r *http.Request) *models.User {
	return r.Context().Value(userKey).(*models.User)
}
