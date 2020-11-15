package server

import (
	"net/http"

	"github.com/asankov/gira/pkg/models"
)

func (s *Server) handleStatusesGet() authorizedHandler {
	return func(w http.ResponseWriter, r *http.Request, user *models.User, token string) {
		s.respond(w, r, models.StatusesResponse{
			Statuses: models.AllStatuses,
		}, http.StatusOK)
	}
}
