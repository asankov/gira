package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSecureHeaders(t *testing.T) {
	srv := &Server{}

	h := srv.secureHeaders(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got, expected := w.Header().Get("X-XSS-Protection"), "1; mode-block"
		if got != expected {
			t.Errorf(`Got ("%s") for "X-XSS-Protection" Header, expected ("%s")`, got, expected)
		}

		got, expected = w.Header().Get("X-Frame-Options"), "deny"
		if got != expected {
			t.Errorf(`Got ("%s") for "X-Frame-Options" Header, expected ("%s")`, got, expected)
		}
	}))

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	h.ServeHTTP(httptest.NewRecorder(), r)
}
