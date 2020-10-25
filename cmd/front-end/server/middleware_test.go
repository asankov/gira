package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	gassert "github.com/asankov/gira/internal/fixtures/assert"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"
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

	assert.True(t, called)
	assert.Equal(t, "deny", w.Header().Get("X-Frame-Options"))
	assert.Equal(t, "1; mode-block", w.Header().Get("X-XSS-Protection"))
}

func TestRequireLogin(t *testing.T) {
	srv := &Server{}

	called := false
	h := srv.requireLogin(authorizedHandler(func(w http.ResponseWriter, r *http.Request, token string) {
		called = true

		assert.Equal(t, token, tokenValue)
	}))

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.AddCookie(&http.Cookie{
		Name:  "token",
		Value: tokenValue,
	})
	h.ServeHTTP(httptest.NewRecorder(), r)

	assert.True(t, called)
}

func TestRequireLoginNoUser(t *testing.T) {
	srv := &Server{}

	called := false
	h := srv.requireLogin(authorizedHandler(func(w http.ResponseWriter, r *http.Request, token string) {
		called = true
	}))

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	h.ServeHTTP(w, r)

	require.False(t, called)
	gassert.Redirect(t, w, "/users/login")
}
