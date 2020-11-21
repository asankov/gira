// +build integration_tests

package main

import (
	"context"
	"testing"

	"github.com/gira-games/client/pkg/client"

	"github.com/stretchr/testify/require"
)

// testCreateAndGetAll creates a user, logs in, then creates two games
// and fetches them.
func testCreateAndGetAll(t *testing.T, cl *client.Client) {
	ctx := context.Background()

	user, err := cl.CreateUser(ctx, &client.CreateUserRequest{
		Email:    "games@test.com",
		Password: "password",
	})
	require.NoError(t, err)

	loginResp, err := cl.LoginUser(ctx, &client.LoginUserRequest{
		Email:    user.Email,
		Password: "password",
	})
	require.NoError(t, err)

	token := loginResp.Token

	batmanGame := createGame(ctx, t, cl, "Batman", token)
	acGame := createGame(ctx, t, cl, "AC", token)

	res, err := cl.GetGames(ctx, &client.GetGamesRequest{Token: token})
	require.NoError(t, err)

	require.Equal(t, 2, len(res.Games))
	require.Contains(t, res.Games, batmanGame)
	require.Contains(t, res.Games, acGame)
}

func createGame(ctx context.Context, t *testing.T, cl *client.Client, name, token string) *client.Game {
	res, err := cl.CreateGame(ctx, &client.CreateGameRequest{
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
