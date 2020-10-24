package server

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/asankov/gira/pkg/models"
	"github.com/asankov/gira/pkg/models/postgres"
	"github.com/gorilla/mux"
	"github.com/hashicorp/go-multierror"
)

type GamePatchRequest struct {
	Status string `json:"status"`
}

var (
	errNameRequired = errors.New("'name' is required parameter")
	errIDNotAllowed = errors.New("'id' is not allowed parameter")
)

func (s *Server) handleGamesCreate() authorizedHandler {
	return func(w http.ResponseWriter, r *http.Request, user *models.User, token string) {
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
			s.Log.Errorf("Error while inserting game into database: %v", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		s.respond(w, r, g, http.StatusOK)
	}
}

func (s *Server) handleGamesGet() authorizedHandler {
	return func(w http.ResponseWriter, r *http.Request, user *models.User, token string) {
		var (
			games []*models.Game
			err   error
		)
		if _, ok := r.URL.Query()["excludeAssigned"]; ok {
			games, err = s.UserGamesModel.GetAvailableGamesFor(user.ID)
		} else {
			games, err = s.GameModel.All()
		}

		if err != nil {
			s.Log.Errorf("Error while fetching games from the database: %v", err)
			http.Error(w, "error fetching games", http.StatusInternalServerError)
			return
		}

		s.respond(w, r, models.GamesResponse{Games: games}, http.StatusOK)
	}
}

func (s *Server) handleGamesGetByID() authorizedHandler {
	return func(w http.ResponseWriter, r *http.Request, user *models.User, token string) {
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
			s.Log.Errorf("Error while fetching game from the database: %v", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		s.respond(w, r, game, http.StatusOK)
	}
}

func validateGame(game *models.Game) error {
	var err *multierror.Error
	if game.Name == "" {
		err = multierror.Append(err, errNameRequired)
	}

	if game.ID != "" {
		err = multierror.Append(err, errIDNotAllowed)
	}

	return err.ErrorOrNil()
}
