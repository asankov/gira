package main

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/asankov/gira/pkg/models"
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

		id, err := s.gameModel.Insert(&game)
		if err != nil {
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		insertedGame, err := s.gameModel.Get(id)
		if err != nil {
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		resp, err := json.Marshal(insertedGame)
		if err != nil {
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
