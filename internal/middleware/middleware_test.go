package middleware

import (
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var logger = log.New(os.Stdout, "", 0)

func TestLog(t *testing.T) {
	called := false
	h := LogRequest(logger)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	}))

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	h.ServeHTTP(w, r)

	got, expected := w.Code, http.StatusOK
	if got != expected {
		t.Errorf(`Got ("%d") for StatusCode, expected ("%d")`, got, expected)
	}

	if !called {
		t.Errorf("Expected next handler to be called, and `called` to be equal to true, instead `called` is false")
	}
}

func TestRecoverPanic(t *testing.T) {

	h := RecoverPanic(logger)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("don't panic")
	}))

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	h.ServeHTTP(w, r)

	got, expected := w.Header().Get("Connection"), "Close"
	if got != expected {
		t.Errorf(`Got ("%s") for "Connection" Header, expected ("%s")`, got, expected)
	}

	if w.Code != http.StatusInternalServerError {
		t.Errorf(`Got ("%d") for StatusCode, expected ("%d")`, w.Code, http.StatusInternalServerError)
	}
}
