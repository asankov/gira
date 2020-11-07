package server

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	gassert "github.com/gira-games/api/internal/fixtures/assert"

	"github.com/gira-games/api/internal/fixtures"
	"github.com/gira-games/api/pkg/models"
	"github.com/golang/mock/gomock"
	"github.com/sirupsen/logrus"
)

var (
	token = "my_test_token"
)

func setupMiddlewareServer(a Authenticator, u UserModel) *Server {
	return &Server{
		Log:           logrus.StandardLogger(),
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
	r.Header.Set(models.XAuthToken, token)

	nextHandlerCalled := false
	h := srv.requireLogin(authorizedHandler(func(w http.ResponseWriter, r *http.Request, u *models.User, token string) {
		nextHandlerCalled = true
	}))

	h.ServeHTTP(w, r)

	gassert.StatusOK(t, w)
	assert.True(t, nextHandlerCalled)
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
				r.Header.Set(models.XAuthToken, token)
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
				r.Header.Set(models.XAuthToken, token)
			},
		},
		{
			name: "Token not present",
			setup: func(a *fixtures.AuthenticatorMock, u *fixtures.UserModelMock, r *http.Request) {
				r.Header.Del(models.XAuthToken)
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
			h := srv.requireLogin(authorizedHandler(func(w http.ResponseWriter, r *http.Request, u *models.User, token string) {
				nextHandlerCalled = true
			}))

			h.ServeHTTP(w, r)

			gassert.StatusCode(t, w, http.StatusUnauthorized)
			assert.False(t, nextHandlerCalled)
		})
	}
}
