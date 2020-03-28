package main

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/asankov/gira/pkg/models"
	"github.com/asankov/gira/pkg/models/postgres"
)

var (
	errNameRequired = errors.New("'name' is required parameter")
	errIDNotAllowed = errors.New("'id' is not allowed parameter")
)

func (s *server) createGameHandler() http.HandlerFunc {
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

		g, err := s.gameModel.Insert(&game)
		if err != nil {
			if errors.Is(err, postgres.ErrNameAlreadyExists) {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			s.log.Printf("error while inserting game into database: %v", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		resp, err := json.Marshal(g)
		if err != nil {
			s.log.Printf("error while encoding response: %v", err)
			http.Error(w, "error encoding response", http.StatusInternalServerError)
			return
		}

		w.Write(resp)
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
