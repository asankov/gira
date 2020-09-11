package assert

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Redirect(t *testing.T, w *httptest.ResponseRecorder, uri string) {
	if w.Code != http.StatusSeeOther {
		t.Errorf("Got (%d) for status code, expected (%d)", w.Code, http.StatusSeeOther)
	}

	got, expected := w.Header().Get("Location"), uri
	if got != expected {
		t.Errorf("Got %s for Location header, expected %s", got, expected)
	}
}

func StatusOK(t *testing.T, w *httptest.ResponseRecorder) {
	StatusCode(t, w, http.StatusOK)
}

func StatusCode(t *testing.T, w *httptest.ResponseRecorder, statusCode int) {
	if w.Code != statusCode {
		t.Errorf("Got (%d) status code, expected (%d)", w.Code, statusCode)
	}
}

func Error(t *testing.T, err error, expectedError error) {
	if err == nil {
		t.Fatal("Got nil error when decoding expired token")
	}
	if !errors.Is(err, expectedError) {
		t.Fatalf("Got (%v) error, expected error to be (%v)", err, expectedError)
	}
}

func NoError(t *testing.T, err error) {
	if err != nil {
		t.Fatalf("Got unexpected error - (%v)", err)
	}
}

func True(t *testing.T, cond bool) {
	if !cond {
		t.Fatalf("Got false, expected true")
	}
}
