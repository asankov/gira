package server

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/asankov/gira/pkg/models"
	"github.com/asankov/gira/pkg/models/postgres"
	"github.com/gorilla/mux"
)

var (
	errNameRequired = errors.New("'name' is required parameter")
	errIDNotAllowed = errors.New("'id' is not allowed parameter")
)

func (s *Server) handleGamesCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var game models.Game

		if err := json.NewDecoder(r.Body).Decode(&game); err != nil {
			http.Error(w, "error decoding body", http.StatusBadRequest)
			return
		}

		if err := validateGame(&game); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		g, err := s.GameModel.Insert(&game)
		if err != nil {
			if errors.Is(err, postgres.ErrNameAlreadyExists) {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			s.Log.Printf("error while inserting game into database: %v", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		s.respond(w, r, g, http.StatusOK)
	}
}

func (s *Server) handleGamesGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		all, err := s.GameModel.All()
		if err != nil {
			s.Log.Printf("error while fetching games from the database: %v", err)
			http.Error(w, "error fetching games", http.StatusInternalServerError)
			return
		}

		s.respond(w, r, models.GamesResponse{Games: all}, http.StatusOK)
	}
}

func (s *Server) handleGamesGetByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		args := mux.Vars(r)
		id := args["id"]
		if id == "" {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		game, err := s.GameModel.Get(id)
		if err != nil {
			if errors.Is(err, postgres.ErrNoRecord) {
				http.NotFound(w, r)
				return
			}
			s.Log.Printf("error while fetching game from the database: %v", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		s.respond(w, r, game, http.StatusOK)
	}
}
func validateGame(game *models.Game) error {
	if game.Name == "" {
		return errNameRequired
	}

	if game.ID != "" {
		return errIDNotAllowed
	}

	return nil
}
