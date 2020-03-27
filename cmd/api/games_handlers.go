package main

import (
	"encoding/json"
	"net/http"

	"github.com/asankov/gira/pkg/models"
)

func (s *server) createGameHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var game models.Game

		if err := json.NewDecoder(r.Body).Decode(&game); err != nil {
			http.Error(w, "error decoding body", http.StatusBadRequest)
			return
		}

		resp, err := json.Marshal(game)
		if err != nil {
			http.Error(w, "error encoding response", http.StatusInternalServerError)
			return
		}

		w.Write(resp)
	}
}
