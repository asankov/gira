package server

import (
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/asankov/gira/internal/fixtures"
	"github.com/golang/mock/gomock"
)

var (
	token = "my_test_token"
)

func setupMiddlewareServer(a Authenticator, u UserModel) *Server {
	return &Server{
		Log:           log.New(os.Stdout, "", 0),
		Authenticator: a,
		UserModel:     u,
	}
}

func TestRequireLogin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	authenticator := fixtures.NewAuthenticatorMock(ctrl)
	userModel := fixtures.NewUserModelMock(ctrl)
	srv := setupMiddlewareServer(authenticator, userModel)

	authenticator.EXPECT().
		DecodeToken(gomock.Eq(token)).
		Return(nil, nil)
	userModel.EXPECT().
		GetUserByToken(gomock.Eq(token)).
		Return(user, nil)

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
		setup func(*fixtures.AuthenticatorMock, *fixtures.UserModelMock, *http.Request)
	}{
		{
			name: "Authenticator error",
			setup: func(a *fixtures.AuthenticatorMock, u *fixtures.UserModelMock, r *http.Request) {
				a.EXPECT().
					DecodeToken(gomock.Eq(token)).
					Return(nil, errors.New("Authenticator error"))
				r.Header.Set("x-auth-token", token)
			},
		},
		{
			name: "DB Error",
			setup: func(a *fixtures.AuthenticatorMock, u *fixtures.UserModelMock, r *http.Request) {
				a.EXPECT().
					DecodeToken(gomock.Eq(token)).
					Return(nil, nil)
				u.EXPECT().
					GetUserByToken(gomock.Eq(token)).
					Return(nil, errors.New("DB Error"))
				r.Header.Set("x-auth-token", token)
			},
		},
		{
			name: "Token not present",
			setup: func(a *fixtures.AuthenticatorMock, u *fixtures.UserModelMock, r *http.Request) {
				r.Header.Del("x-auth-token")
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			authenticator := fixtures.NewAuthenticatorMock(ctrl)
			userModel := fixtures.NewUserModelMock(ctrl)
			srv := setupMiddlewareServer(authenticator, userModel)

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/", nil)

			testCase.setup(authenticator, userModel, r)

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
