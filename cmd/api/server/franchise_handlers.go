package server

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gira-games/api/pkg/models"
	"github.com/gira-games/api/pkg/models/postgres"
	"github.com/hashicorp/go-multierror"
)

func (s *Server) handleFranchisesGet() authorizedHandler {
	return func(w http.ResponseWriter, r *http.Request, user *models.User, token string) {
		franchises, err := s.FranchiseModel.All()

		if err != nil {
			s.Log.Errorf("Error while fetching franchises from the database: %v", err)
			s.internalError(w, r)
			return
		}

		s.respond(w, r, models.FranchisesResponse{Franchises: franchises}, http.StatusOK)
	}
}

func (s *Server) handleFranchisesCreate() authorizedHandler {
	return func(w http.ResponseWriter, r *http.Request, user *models.User, token string) {
		var franchise models.Franchise

		if err := json.NewDecoder(r.Body).Decode(&franchise); err != nil {
			s.respondError(w, r, "Error decoding body", http.StatusBadRequest)
			return
		}

		if err := validateFranchise(&franchise); err != nil {
			s.respondError(w, r, err.Error(), http.StatusBadRequest)
			return
		}

		g, err := s.FranchiseModel.Insert(&franchise)
		if err != nil {
			if errors.Is(err, postgres.ErrNameAlreadyExists) {
				s.respondError(w, r, "Franchise with the same name already exists", http.StatusBadRequest)
				return
			}
			s.Log.Errorf("Error while inserting game into database: %v", err)
			s.internalError(w, r)
			return
		}

		s.respond(w, r, g, http.StatusOK)
	}
}

func validateFranchise(franchise *models.Franchise) error {
	var err *multierror.Error
	if franchise.ID != "" {
		err = multierror.Append(err, errIDNotAllowed)
	}
	if franchise.Name == "" {
		err = multierror.Append(err, errNameRequired)
	}

	return err.ErrorOrNil()
}
