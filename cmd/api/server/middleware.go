package server

import (
	"net/http"

	"github.com/gira-games/api/pkg/models"
)

type authorizedHandler func(http.ResponseWriter, *http.Request, *models.User, string)

func (s *Server) requireLogin(next authorizedHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get(models.XAuthToken)
		if token == "" {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		if _, err := s.Authenticator.DecodeToken(token); err != nil {
			s.Log.Errorf("Error while decoding token: %v", err)
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		// this will fail if the token is invalidated
		user, err := s.UserModel.GetUserByToken(token)
		if err != nil {
			// TODO: if err == not found else ..
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		next(w, r, user, token)
	})
}
