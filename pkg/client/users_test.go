package client_test

import (
	"net/http"
	"testing"

	"github.com/gira-games/api/pkg/models"

	"github.com/gira-games/api/internal/fixtures"
	"github.com/stretchr/testify/require"
)

var (
	user = &models.User{
		ID:       "123",
		Username: "test_user",
		Email:    "test@mail.com",
	}
	userResponse = &models.UserResponse{
		User: user,
	}
	userLoginResponse = &models.UserLoginResponse{
		Token: token,
	}
)

func TestGetUser(t *testing.T) {
	ts := fixtures.NewTestServer(t).
		Path("/users").
		Token(token).
		Data(userResponse).
		Build()
	defer ts.Close()

	cl := newClient(t, ts.URL)

	usr, err := cl.GetUser(token)
	require.NoError(t, err)
	require.Equal(t, usr, user)
}

func TestCreateUser(t *testing.T) {
	ts := fixtures.NewTestServer(t).
		Path("/users").
		Method(http.MethodPost).
		Data(user).
		Build()
	defer ts.Close()

	cl := newClient(t, ts.URL)

	usr, err := cl.CreateUser(user)
	require.NoError(t, err)
	require.Equal(t, usr, user)
}

func TestLogin(t *testing.T) {
	ts := fixtures.NewTestServer(t).
		Path("/users/login").
		Method(http.MethodPost).
		Data(userLoginResponse).
		Build()
	defer ts.Close()

	cl := newClient(t, ts.URL)

	usrLoginResponse, err := cl.LoginUser(user)
	require.NoError(t, err)
	require.Equal(t, usrLoginResponse, userLoginResponse)
}

func TestLogout(t *testing.T) {
	ts := fixtures.NewTestServer(t).
		Method(http.MethodPost).
		Path("/users/logout").
		Token(token).
		Build()
	defer ts.Close()

	cl := newClient(t, ts.URL)

	err := cl.LogoutUser(token)
	require.NoError(t, err)
}
