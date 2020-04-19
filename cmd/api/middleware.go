package main

import "net/http"

func (s *server) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: better format for the request
		s.log.Printf("Incoming request: %#v\n", r)

		next.ServeHTTP(w, r)
	})
}
