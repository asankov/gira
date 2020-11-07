package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	gassert "github.com/gira-games/api/internal/fixtures/assert"
	"github.com/stretchr/testify/assert"

	"github.com/sirupsen/logrus"
)

var logger = logrus.StandardLogger()

func TestLog(t *testing.T) {
	called := false
	h := LogRequest(logger)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	}))

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	h.ServeHTTP(w, r)

	gassert.StatusOK(t, w)
	assert.True(t, called)
}

func TestRecoverPanic(t *testing.T) {

	h := RecoverPanic(logger)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("don't panic")
	}))

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	h.ServeHTTP(w, r)

	assert.Equal(t, "Close", w.Header().Get("Connection"))
	gassert.StatusCode(t, w, http.StatusInternalServerError)
}
