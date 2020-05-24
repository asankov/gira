package middleware

import (
	"log"
	"net/http"
)

// LogRequest returns a function that accepts an http.Handler and returns another http.Handler
// that logs all requests that come in.
// This format is needed so that we can pass an external logger and reuse this function
// in an alice middleware chain.
func LogRequest(log *log.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Printf("%s - %s %s %s\n", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())

			next.ServeHTTP(w, r)
		})
	}
}
