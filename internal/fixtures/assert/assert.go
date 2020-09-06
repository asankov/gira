package assert

import (
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
