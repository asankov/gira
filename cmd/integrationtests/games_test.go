// +build integration_tests

package main

import (
	"context"
	"testing"

	"github.com/asankov/gira/pkg/client"

	"github.com/stretchr/testify/require"
)

// testCreateAndGetAll creates a user, logs in, then creates two games
// and fetches them.
func testCreateAndGetAll(t *testing.T, cl *client.Client) {
	user, err := cl.CreateUser(context.Background(), &client.CreateUserRequest{
		Email:    "games@test.com",
		Password: "password",
	})
	require.NoError(t, err)

	loginResp, err := cl.LoginUser(context.Background(), &client.LoginUserRequest{
		Email:    user.Email,
		Password: "password",
	})
	require.NoError(t, err)

	token := loginResp.Token

	batmanGame := createGame(t, cl, "Batman", token)
	acGame := createGame(t, cl, "AC", token)

	res, err := cl.GetGames(context.Background(), &client.GetGamesRequest{Token: token})
	require.NoError(t, err)

	require.Equal(t, 2, len(res.Games))
	require.Contains(t, res.Games, batmanGame)
	require.Contains(t, res.Games, acGame)
}

func createGame(t *testing.T, cl *client.Client, name, token string) *client.Game {
	res, err := cl.CreateGame(context.Background(), &client.CreateGameRequest{
		Token: token,
		Game: &client.Game{
			Name: name,
		},
	})
	require.NoError(t, err)
	require.NotEmpty(t, res.Game.ID)
	require.Empty(t, res.Game.FranchiseID)
	require.Equal(t, name, res.Game.Name)

	return res.Game
}
