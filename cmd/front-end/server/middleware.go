package server

import "net/http"

// TODO: this whole file is copied from cmd/api/middleware.go
// find a way to refactor it and reduce the duplication

func (s *Server) secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-XSS-Protection", "1; mode-block")
		w.Header().Set("X-Frame-Options", "deny")

		next.ServeHTTP(w, r)
	})
}

func (s *Server) requireLogin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := r.Cookie("token"); err != nil {
			w.Header().Add("Location", "/users/login")
			w.WriteHeader(http.StatusSeeOther)
			return
		}

		// at this point we don't care whether the cookie is valid or not, just that is exists
		// if the token inside the cookie is not valid the back-end would return 401 Unathorized

		next.ServeHTTP(w, r)
	})
}
