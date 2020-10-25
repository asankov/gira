// +build integration_tests

package integrationtests

import (
	"testing"

	"github.com/asankov/gira/pkg/models"
	"github.com/stretchr/testify/require"
)

// TestUserLifecycle test the lifecycle of the user.
// It creates a user, logs in, gets user info via the token, received on login,
// logs out, and then checks that after logging out the token has been invalidated.
func TestUserLifecycle(t *testing.T) {
	cl := setup(t)

	user, err := cl.CreateUser(&models.User{
		Email:    "integration@test.com",
		Password: "pass",
	})
	require.NoError(t, err)
	require.Equal(t, "integration@test.com", user.Email)
	require.Equal(t, "integration@test.com", user.Username, "the username should be the same as the email by default")
	require.Empty(t, user.Password, "the server should not return the user password")

	resp, err := cl.LoginUser(&models.User{
		Email:    "integration@test.com",
		Password: "pass",
	})
	require.NoError(t, err)
	require.NotEmpty(t, resp.Token)

	user, err = cl.GetUser(resp.Token)
	require.NoError(t, err)
	require.Equal(t, "integration@test.com", user.Email)
	require.Equal(t, "integration@test.com", user.Username, "the username should be the same as the email by default")
	require.Empty(t, user.Password, "the server should not return the user password")

	err = cl.LogoutUser(resp.Token)
	require.NoError(t, err)

	user, err = cl.GetUser(resp.Token)
	require.Nil(t, user)
	// TODO: assert error once we start returning proper errors
	require.Error(t, err)
}
