package server

import (
	"context"
	"net/http"
)

type contextUserKeyType string

var contextUserKey contextUserKeyType

func (s *Server) requireLogin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("x-auth-token")
		if token == "" {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		user, err := s.Authenticator.DecodeToken(token)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), contextUserKey, user)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
