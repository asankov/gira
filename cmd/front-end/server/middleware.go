package server

import (
	"context"
	"net/http"
)

type contextTokenKeyType string

var (
	contextTokenKey contextTokenKeyType
)

func (s *Server) secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-XSS-Protection", "1; mode-block")
		w.Header().Set("X-Frame-Options", "deny")

		next.ServeHTTP(w, r)
	})
}

func (s *Server) requireLogin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := r.Cookie("token")
		if err != nil {
			w.Header().Add("Location", "/users/login")
			w.WriteHeader(http.StatusSeeOther)
			return
		}

		ctx := context.WithValue(r.Context(), contextTokenKey, token.Value)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
