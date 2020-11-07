package server

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	gassert "github.com/gira-games/api/internal/fixtures/assert"
	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"

	"github.com/gira-games/api/internal/auth"
	"github.com/gira-games/api/internal/fixtures"
	"github.com/gira-games/api/pkg/models"
	"github.com/gira-games/api/pkg/models/postgres"
	"github.com/golang/mock/gomock"
	"github.com/sirupsen/logrus"
)

var (
	expectedUser = models.User{
		Username: "test",
		Email:    "test@test.com",
		Password: "t3$T123",
	}
)

func newServer(t *testing.T, opts *Options) *Server {
	if opts == nil {
		opts = &Options{}
	}
	opts.Log = logrus.StandardLogger()

	srv, err := New(opts)
	require.Nil(t, err)

	return srv
}

func TestUserCreate(t *testing.T) {
	testCases := []struct {
		Name          string
		UserInRequest models.User
		ExpectedUser  models.User
	}{
		{
			Name:          "Pass username, email and password",
			UserInRequest: expectedUser,
			ExpectedUser:  expectedUser,
		},
		{
			Name: "Pass enail and password, username gets filled automatically",
			UserInRequest: models.User{
				Email:    "test@mail.com",
				Password: "pass",
			},
			ExpectedUser: models.User{
				Username: "test@mail.com",
				Email:    "test@mail.com",
				Password: "pass",
			},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userModel := fixtures.NewUserModelMock(ctrl)
			authenticator := fixtures.NewAuthenticatorMock(ctrl)
			srv := newServer(t, &Options{
				UserModel:     userModel,
				Authenticator: authenticator,
			})

			userModel.EXPECT().
				Insert(&testCase.ExpectedUser).
				Return(&testCase.ExpectedUser, nil)

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, "/users", fixtures.Marshal(t, testCase.UserInRequest))
			srv.ServeHTTP(w, r)

			gassert.StatusOK(t, w)

			var user models.User
			fixtures.Decode(t, w.Body, &user)

			assert.Equal(t, user.Email, testCase.ExpectedUser.Email)
			assert.Equal(t, user.Username, testCase.ExpectedUser.Username)
		})
	}
}

func TestUserCreateValidationError(t *testing.T) {
	cases := []struct {
		name string
		user *models.User
	}{
		{
			name: "No email",
			user: &models.User{
				Username: "test",
				Password: "t3$t",
			},
		},
		{
			name: "No password",
			user: &models.User{
				Username: "test",
				Email:    "test@test.com",
			},
		},
		{
			name: "Filled ID",
			user: &models.User{
				ID:       "1",
				Username: "test",
				Email:    "test@test.com",
				Password: "t3$t",
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			srv := newServer(t, nil)

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, "/users", fixtures.Marshal(t, c.user))
			srv.ServeHTTP(w, r)

			gassert.StatusCode(t, w, http.StatusBadRequest)
		})
	}
}

func TestUserCreateEmptyBody(t *testing.T) {
	srv := newServer(t, nil)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/users", nil)
	srv.ServeHTTP(w, r)

	gassert.StatusCode(t, w, http.StatusBadRequest)
}

func TestUserCreateDBError(t *testing.T) {
	testCases := []struct {
		name         string
		dbError      error
		expectedCode int
	}{
		{
			name:         "Email already exists",
			dbError:      postgres.ErrEmailAlreadyExists,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "Name already exists",
			dbError:      postgres.ErrUsernameAlreadyExists,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "Unknown error",
			dbError:      errors.New("unknown error"),
			expectedCode: http.StatusInternalServerError,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userModel := fixtures.NewUserModelMock(ctrl)

			srv := newServer(t, &Options{
				UserModel: userModel,
			})

			userModel.EXPECT().
				Insert(&expectedUser).
				Return(nil, testCase.dbError)

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, "/users", fixtures.Marshal(t, expectedUser))
			srv.ServeHTTP(w, r)

			gassert.StatusCode(t, w, testCase.expectedCode)

		})
	}
}

func TestUserLogin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userModel := fixtures.NewUserModelMock(ctrl)
	authenticatorMock := fixtures.NewAuthenticatorMock(ctrl)

	srv := newServer(t, &Options{
		UserModel:     userModel,
		Authenticator: authenticatorMock,
	})

	userModel.EXPECT().
		Authenticate(expectedUser.Email, expectedUser.Password).
		Return(&expectedUser, nil)
	userModel.EXPECT().
		AssociateTokenWithUser(expectedUser.ID, token).
		Return(nil)

	token := "my_test_token"
	authenticatorMock.EXPECT().
		NewTokenForUser(&expectedUser).
		Return(token, nil)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/users/login", fixtures.Marshal(t, expectedUser))
	srv.ServeHTTP(w, r)

	gassert.StatusOK(t, w)
	var userResponse models.UserLoginResponse
	fixtures.Decode(t, w.Body, &userResponse)
	assert.Equal(t, token, userResponse.Token)
}

func TestUserLoginValidationError(t *testing.T) {
	testCases := []struct {
		name string
		user *models.User
	}{
		{
			name: "No email",
			user: &models.User{
				Email:    "",
				Password: "T3$T",
			},
		},
		{
			name: "No password",
			user: &models.User{
				Email:    "test@mail.com",
				Password: "",
			},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userModel := fixtures.NewUserModelMock(ctrl)

			srv := newServer(t, &Options{
				UserModel: userModel,
			})

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, "/users/login", fixtures.Marshal(t, testCase.user))
			srv.ServeHTTP(w, r)

			gassert.StatusCode(t, w, http.StatusBadRequest)

			// TODO: assert body, once we start returning proper errors
		})
	}
}

func TestUserLoginServiceError(t *testing.T) {
	testCases := []struct {
		name         string
		setup        func(u *fixtures.UserModelMock, a *fixtures.AuthenticatorMock)
		expectedCode int
	}{
		{
			name: "UserModel.Authenticate fails",
			setup: func(u *fixtures.UserModelMock, a *fixtures.AuthenticatorMock) {
				u.EXPECT().
					Authenticate(expectedUser.Email, expectedUser.Password).
					Return(nil, errors.New("user not found"))
			},
			expectedCode: http.StatusUnauthorized,
		},
		{
			name: "Authenticator.NewTokenForUser fails",
			setup: func(u *fixtures.UserModelMock, a *fixtures.AuthenticatorMock) {
				u.EXPECT().
					Authenticate(expectedUser.Email, expectedUser.Password).
					Return(&expectedUser, nil)

				a.EXPECT().
					NewTokenForUser(&expectedUser).
					Return("", errors.New("intentional error"))
			},
			expectedCode: http.StatusInternalServerError,
		},
		{
			name: "UserModel.AssociateTokenWithUser fails",
			setup: func(u *fixtures.UserModelMock, a *fixtures.AuthenticatorMock) {
				u.EXPECT().
					Authenticate(expectedUser.Email, expectedUser.Password).
					Return(&expectedUser, nil)

				a.EXPECT().
					NewTokenForUser(&expectedUser).
					Return(token, nil)

				u.EXPECT().
					AssociateTokenWithUser(expectedUser.ID, token).
					Return(errors.New("intentional error while associating token with user"))
			},
			expectedCode: http.StatusInternalServerError,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userModel := fixtures.NewUserModelMock(ctrl)
			authenticatorMock := fixtures.NewAuthenticatorMock(ctrl)

			testCase.setup(userModel, authenticatorMock)

			srv := newServer(t, &Options{
				UserModel:     userModel,
				Authenticator: authenticatorMock,
			})

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, "/users/login", fixtures.Marshal(t, expectedUser))
			srv.ServeHTTP(w, r)

			gassert.StatusCode(t, w, testCase.expectedCode)
		})
	}
}

func TestUserGet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userModelMock := fixtures.NewUserModelMock(ctrl)
	authenticatorMock := fixtures.NewAuthenticatorMock(ctrl)

	authenticatorMock.EXPECT().
		DecodeToken(gomock.Eq(token)).
		Return(nil, nil)
	userModelMock.EXPECT().
		GetUserByToken(gomock.Eq(token)).
		Return(&expectedUser, nil)

	srv := newServer(t, &Options{
		UserModel:     userModelMock,
		Authenticator: authenticatorMock,
	})

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/users", nil)
	r.Header.Add(models.XAuthToken, token)
	srv.ServeHTTP(w, r)

	gassert.StatusOK(t, w)

	var userResponse models.UserResponse
	fixtures.Decode(t, w.Body, &userResponse)
	gotUser := userResponse.User

	assert.Equal(t, gotUser.ID, expectedUser.ID)
	assert.Equal(t, gotUser.Email, expectedUser.Email)
	assert.Equal(t, gotUser.Username, expectedUser.Username)
}

func TestUserGetUnathorized(t *testing.T) {
	testCases := []struct {
		name  string
		setup func(*fixtures.AuthenticatorMock, *fixtures.UserModelMock, *http.Request)
	}{
		{
			name:  "No token",
			setup: func(a *fixtures.AuthenticatorMock, u *fixtures.UserModelMock, r *http.Request) {},
		},
		{
			name: "Invalid signature",
			setup: func(a *fixtures.AuthenticatorMock, u *fixtures.UserModelMock, r *http.Request) {
				r.Header.Add(models.XAuthToken, token)

				a.EXPECT().
					DecodeToken(gomock.Eq(token)).
					Return(nil, auth.ErrInvalidSignature)
			},
		},
		{
			name: "Token expired",
			setup: func(a *fixtures.AuthenticatorMock, u *fixtures.UserModelMock, r *http.Request) {
				r.Header.Add(models.XAuthToken, token)

				a.EXPECT().
					DecodeToken(gomock.Eq(token)).
					Return(nil, auth.ErrTokenExpired)
			},
		},
		{
			name: "Generic token error",
			setup: func(a *fixtures.AuthenticatorMock, u *fixtures.UserModelMock, r *http.Request) {
				r.Header.Add(models.XAuthToken, token)

				a.EXPECT().
					DecodeToken(gomock.Eq(token)).
					Return(nil, errors.New("some generic error"))
			},
		},
		{
			name: "UserModel error",
			setup: func(a *fixtures.AuthenticatorMock, u *fixtures.UserModelMock, r *http.Request) {
				r.Header.Add(models.XAuthToken, token)

				a.EXPECT().
					DecodeToken(gomock.Eq(token)).
					Return(nil, nil)
				u.EXPECT().
					GetUserByToken(gomock.Eq(token)).
					Return(nil, errors.New("some error"))
			},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			authenticatorMock := fixtures.NewAuthenticatorMock(ctrl)
			userModelMock := fixtures.NewUserModelMock(ctrl)
			srv := newServer(t, &Options{
				Authenticator: authenticatorMock,
				UserModel:     userModelMock,
			})

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/users", nil)
			testCase.setup(authenticatorMock, userModelMock, r)
			srv.ServeHTTP(w, r)

			gassert.StatusCode(t, w, http.StatusUnauthorized)

			// TODO: assert error once we return JSON
		})
	}
}

func TestUserLogout(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userModelMock := fixtures.NewUserModelMock(ctrl)
	authenticatorMock := fixtures.NewAuthenticatorMock(ctrl)
	srv := newServer(t, &Options{
		UserModel:     userModelMock,
		Authenticator: authenticatorMock,
	})

	authenticatorMock.EXPECT().
		DecodeToken(gomock.Eq(token)).
		Return(nil, nil)
	userModelMock.EXPECT().
		GetUserByToken(gomock.Eq(token)).
		Return(&expectedUser, nil)
	userModelMock.EXPECT().
		InvalidateToken(gomock.Eq(expectedUser.ID), gomock.Eq(token)).
		Return(nil)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/users/logout", nil)
	r.Header.Add(models.XAuthToken, token)
	srv.ServeHTTP(w, r)

	gassert.StatusOK(t, w)
}

func TestUserLogoutInvalidateError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userModelMock := fixtures.NewUserModelMock(ctrl)
	authenticatorMock := fixtures.NewAuthenticatorMock(ctrl)
	srv := newServer(t, &Options{
		UserModel:     userModelMock,
		Authenticator: authenticatorMock,
	})

	authenticatorMock.EXPECT().
		DecodeToken(gomock.Eq(token)).
		Return(nil, nil)
	userModelMock.EXPECT().
		GetUserByToken(gomock.Eq(token)).
		Return(&expectedUser, nil)
	userModelMock.EXPECT().
		InvalidateToken(gomock.Eq(expectedUser.ID), gomock.Eq(token)).
		Return(errors.New("this token is already validated"))

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/users/logout", nil)
	r.Header.Add(models.XAuthToken, token)
	srv.ServeHTTP(w, r)

	gassert.StatusCode(t, w, http.StatusBadRequest)
}
