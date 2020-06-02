package server

import (
	"context"
	"net/http"
)

type contextUserKeyType string
type contextTokenKeyType string

var contextUserKey contextUserKeyType
var contextTokenKey contextTokenKeyType

func (s *Server) requireLogin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("x-auth-token")
		if token == "" {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		if _, err := s.Authenticator.DecodeToken(token); err != nil {
			s.Log.Printf("error while decoding token: %v", err)
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

		ctx := context.WithValue(r.Context(), contextUserKey, user)
		ctx = context.WithValue(ctx, contextTokenKey, token)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
