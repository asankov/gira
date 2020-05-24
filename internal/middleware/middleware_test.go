package middleware

import (
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestRecoverPanic(t *testing.T) {

	h := RecoverPanic(log.New(os.Stdout, "", 0))(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
