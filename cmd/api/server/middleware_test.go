package server

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/asankov/gira/internal/fixtures"
	"github.com/golang/mock/gomock"
)

var (
	token = "my_test_token"
)

func TestSecureHeaders(t *testing.T) {
	srv := newServer(nil, nil, nil)

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

func TestRecoverPanic(t *testing.T) {
	// TODO: mock logger and assert output, once server.Log is made an interface
	srv := newServer(nil, nil, nil)

	h := srv.recoverPanic(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

func TestRequireLogin(t *testing.T) {
	ctrl := gomock.NewController(t)
	authenticator := fixtures.NewAuthenticatorMock(ctrl)
	srv := newServer(nil, nil, authenticator)

	authenticator.EXPECT().
		DecodeToken(gomock.Eq(token)).
		Return(nil, nil)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set("x-auth-token", token)

	nextHandlerCalled := false
	h := srv.requireLogin(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextHandlerCalled = true
	}))

	h.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Errorf(`Got ("%d") for StatusCode, expected ("%d")`, w.Code, http.StatusOK)
	}
	if !nextHandlerCalled {
		t.Errorf("Got false for nextHandlerCalled, expected true")
	}
}
func TestRequireLoginError(t *testing.T) {
	testCases := []struct {
		name  string
		setup func(a *fixtures.AuthenticatorMock, r *http.Request)
	}{
		{
			name: "Authenticator error",
			setup: func(a *fixtures.AuthenticatorMock, r *http.Request) {
				a.EXPECT().
					DecodeToken(gomock.Eq(token)).
					Return(nil, errors.New("Authenticator error"))
				r.Header.Set("x-auth-token", token)
			},
		},
		{
			name: "Token not present",
			setup: func(a *fixtures.AuthenticatorMock, r *http.Request) {
				r.Header.Del("x-auth-token")
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			authenticator := fixtures.NewAuthenticatorMock(ctrl)
			srv := newServer(nil, nil, authenticator)

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/", nil)

			testCase.setup(authenticator, r)

			nextHandlerCalled := false
			h := srv.requireLogin(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				nextHandlerCalled = true
			}))

			h.ServeHTTP(w, r)

			if w.Code != http.StatusUnauthorized {
				t.Errorf(`Got ("%d") for StatusCode, expected ("%d")`, w.Code, http.StatusUnauthorized)
			}
			if nextHandlerCalled {
				t.Errorf("Got true for nextHandlerCalled, expected false")
			}
		})
	}
}
