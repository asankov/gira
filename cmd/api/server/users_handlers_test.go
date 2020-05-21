package server

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/asankov/gira/internal/auth"
	"github.com/asankov/gira/internal/fixtures"
	"github.com/asankov/gira/pkg/models"
	"github.com/asankov/gira/pkg/models/postgres"
	"github.com/golang/mock/gomock"
)

var (
	expectedUser = models.User{
		Username: "test",
		Email:    "test@test.com",
		Password: "t3$T123",
	}
)

func TestUserCreate(t *testing.T) {
	ctrl := gomock.NewController(t)

	userModel := fixtures.NewUserModelMock(ctrl)

	srv := newServer(nil, userModel, auth.NewAutheniticator("some_secret"))

	userModel.EXPECT().
		Insert(&expectedUser).
		Return(&expectedUser, nil)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/users", marshall(t, expectedUser))
	srv.ServeHTTP(w, r)

	got, expected := w.Code, http.StatusOK
	if got != expected {
		t.Fatalf("Got (%d) for status code, expected (%d)", got, expected)
	}

	var user models.User
	decode(t, w, &user)
	if user.Username != expectedUser.Username {
		t.Errorf("Got (%s) for username, expected (%s)", user.Username, expectedUser.Username)
	}
	if user.Email != expectedUser.Email {
		t.Errorf("Got (%s) for email, expected (%s)", user.Email, expectedUser.Email)
	}
}

func TestUserCreateValidationError(t *testing.T) {
	cases := []struct {
		name string
		user *models.User
	}{
		{
			name: "No username",
			user: &models.User{
				Email:    "test@test.com",
				Password: "t3$t",
			},
		},
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
		{
			name: "Filled hashed password",
			user: &models.User{
				Username:       "test",
				Email:          "test@test.com",
				Password:       "t3$t",
				HashedPassword: []byte("t3$t"),
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			srv := newServer(nil, nil, nil)

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, "/users", marshall(t, c.user))
			srv.ServeHTTP(w, r)

			got, expected := w.Code, http.StatusBadRequest
			if got != expected {
				t.Fatalf("Got (%d) for status code, expected (%d)", got, expected)
			}
		})
	}
}

func TestUserCreateDBError(t *testing.T) {
	cases := []struct {
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
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			userModel := fixtures.NewUserModelMock(ctrl)

			srv := newServer(nil, userModel, nil)

			userModel.EXPECT().
				Insert(&expectedUser).
				Return(nil, c.dbError)

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, "/users", marshall(t, expectedUser))
			srv.ServeHTTP(w, r)

			got, expected := w.Code, c.expectedCode
			if got != expected {
				t.Fatalf("Got (%d) for status code, expected (%d)", got, expected)
			}
		})
	}
}
