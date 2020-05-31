package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

var tokenValue = "my_token"

func TestSecureHeaders(t *testing.T) {
	srv := &Server{}

	called := false
	h := srv.secureHeaders(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	}))

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	h.ServeHTTP(w, r)

	if !called {
		t.Errorf("Expected next handler to be called, and `called` to be equal to true, instead `called` is false")
	}

	got, expected := w.Header().Get("X-XSS-Protection"), "1; mode-block"
	if got != expected {
		t.Errorf(`Got ("%s") for "X-XSS-Protection" Header, expected ("%s")`, got, expected)
	}

	got, expected = w.Header().Get("X-Frame-Options"), "deny"
	if got != expected {
		t.Errorf(`Got ("%s") for "X-Frame-Options" Header, expected ("%s")`, got, expected)
	}
}

func TestRequireLogin(t *testing.T) {
	srv := &Server{}

	called := false
	h := srv.requireLogin(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true

		token, ok := r.Context().Value(contextTokenKey).(string)
		if !ok {
			t.Errorf("Expected `r.Context().Value(contextTokenKey)` to be of type string")
		}
		if token != tokenValue {
			t.Errorf("Got (%s) for token from context, expected (%s)", token, tokenValue)
		}
	}))

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.AddCookie(&http.Cookie{
		Name:  "token",
		Value: tokenValue,
	})
	h.ServeHTTP(httptest.NewRecorder(), r)

	if !called {
		t.Errorf("Expected next handler to be called, and `called` to be equal to true, instead `called` is false")
	}
}

func TestRequireLoginNoUser(t *testing.T) {
	srv := &Server{}

	called := false
	h := srv.requireLogin(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	}))

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	h.ServeHTTP(w, r)

	if called {
		t.Errorf("Expected next handler to not be called, and `called` to be equal to false, instead `called` is true")
	}

	got, expected := w.Code, http.StatusSeeOther
	if got != expected {
		t.Errorf("Got status code (%d), expected (%d)", got, expected)
	}

	gotH, expectedH := w.Header().Get("Location"), "/users/login"
	if gotH != expectedH {
		t.Errorf("Got (%s) for Location header, expected (%s)", gotH, expectedH)
	}
}
