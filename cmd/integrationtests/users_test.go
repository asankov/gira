// +build integration_tests

package main

import (
	"context"
	"testing"

	"github.com/asankov/gira/pkg/client"

	"github.com/stretchr/testify/require"
)

// testUserLifecycle test the lifecycle of the user.
// It creates a user, logs in, gets user info via the token, received on login,
// logs out, and then checks that after logging out the token has been invalidated.
func testUserLifecycle(t *testing.T, cl *client.Client) {
	user, err := cl.CreateUser(context.Background(), &client.CreateUserRequest{
		Email:    "integration@test.com",
		Password: "pass",
	})
	require.NoError(t, err)
	require.Equal(t, "integration@test.com", user.Email)
	require.Equal(t, "integration@test.com", user.Username, "the username should be the same as the email by default")

	resp, err := cl.LoginUser(context.Background(), &client.LoginUserRequest{
		Email:    "integration@test.com",
		Password: "pass",
	})
	require.NoError(t, err)
	require.NotEmpty(t, resp.Token)

	getUserResp, err := cl.GetUser(context.Background(), &client.GetUserRequest{Token: resp.Token})
	require.NoError(t, err)
	require.Equal(t, "integration@test.com", getUserResp.Email)
	require.Equal(t, "integration@test.com", getUserResp.Username, "the username should be the same as the email by default")

	err = cl.LogoutUser(context.Background(), &client.LogoutUserRequest{Token: resp.Token})
	require.NoError(t, err)

	getUserResp, err = cl.GetUser(context.Background(), &client.GetUserRequest{Token: resp.Token})
	require.Nil(t, getUserResp)
	// TODO: assert error once we start returning proper errors
	require.Error(t, err)
}
