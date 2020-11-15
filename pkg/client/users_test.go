package client_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/asankov/gira/internal/fixtures"
	"github.com/asankov/gira/pkg/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	user = &client.User{
		ID:       "123",
		Username: "test_user",
		Email:    "test@mail.com",
	}
	userResponse = struct {
		User *client.User `json:"user"`
	}{
		User: user,
	}
	userLoginResponse = &client.UserLoginResponse{
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

	resp, err := cl.GetUser(context.Background(), &client.GetUserRequest{Token: token})
	require.NoError(t, err)
	assert.Equal(t, resp.Email, user.Email)
	assert.Equal(t, resp.ID, user.ID)
	assert.Equal(t, resp.Username, user.Username)
}

func TestCreateUser(t *testing.T) {
	ts := fixtures.NewTestServer(t).
		Path("/users").
		Method(http.MethodPost).
		Data(user).
		Build()
	defer ts.Close()

	cl := newClient(t, ts.URL)

	resp, err := cl.CreateUser(context.Background(), &client.CreateUserRequest{
		Username: user.Username,
		Email:    user.Email,
		Password: user.Password,
	})
	require.NoError(t, err)
	assert.Equal(t, resp.Email, user.Email)
	assert.Equal(t, resp.ID, user.ID)
	assert.Equal(t, resp.Username, user.Username)
}

func TestLogin(t *testing.T) {
	ts := fixtures.NewTestServer(t).
		Path("/users/login").
		Method(http.MethodPost).
		Data(userLoginResponse).
		Build()
	defer ts.Close()

	cl := newClient(t, ts.URL)

	resp, err := cl.LoginUser(context.Background(), &client.LoginUserRequest{
		Username: user.Username,
		Email:    user.Email,
		Password: user.Password,
	})
	require.NoError(t, err)
	require.Equal(t, resp, userLoginResponse)
}

func TestLogout(t *testing.T) {
	ts := fixtures.NewTestServer(t).
		Method(http.MethodPost).
		Path("/users/logout").
		Token(token).
		Build()
	defer ts.Close()

	cl := newClient(t, ts.URL)

	err := cl.LogoutUser(context.Background(), &client.LogoutUserRequest{
		Token: token,
	})
	require.NoError(t, err)
}
