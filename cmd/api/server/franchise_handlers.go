package server

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/asankov/gira/pkg/models"
	"github.com/asankov/gira/pkg/models/postgres"
	"github.com/hashicorp/go-multierror"
)

func (s *Server) handleFranchisesGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		franchises, err := s.FranchiseModel.All()

		if err != nil {
			s.Log.Errorf("Error while fetching franchises from the database: %v", err)
			// TODO: json
			http.Error(w, "error fetching franchises", http.StatusInternalServerError)
			return
		}

		s.respond(w, r, models.FranchisesResponse{Franchises: franchises}, http.StatusOK)
	}
}

func (s *Server) handleFranchisesCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var franchise models.Franchise

		if err := json.NewDecoder(r.Body).Decode(&franchise); err != nil {
			http.Error(w, "error decoding body", http.StatusBadRequest)
			return
		}

		if err := validateFranchise(&franchise); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		g, err := s.FranchiseModel.Insert(&franchise)
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

func validateFranchise(franchise *models.Franchise) error {
	var err *multierror.Error
	if franchise.Name == "" {
		err = multierror.Append(err, errNameRequired)
	}

	return err.ErrorOrNil()
}
