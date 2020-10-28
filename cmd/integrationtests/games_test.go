// +build integration_tests

package integrationtests

import (
	"testing"

	"github.com/asankov/gira/pkg/client"
	"github.com/asankov/gira/pkg/models"
	"github.com/stretchr/testify/require"
)

// TestCreateAndGetAll creates a user, logs in, then creates two games
// and fetches them.
func TestCreateAndGetAll(t *testing.T) {
	cl := setup(t)

	user, err := cl.CreateUser(&models.User{
		Email:    "games@test.com",
		Password: "password",
	})
	require.NoError(t, err)

	loginResp, err := cl.LoginUser(&models.User{
		Email:    user.Email,
		Password: "password",
	})
	require.NoError(t, err)

	token := loginResp.Token

	batmanGame := createGame(t, cl, "Batman", token)
	acGame := createGame(t, cl, "AC", token)

	games, err := cl.GetGames(token, client.NoOptions)
	require.NoError(t, err)

	require.Equal(t, 2, len(games))
	require.Contains(t, games, batmanGame)
	require.Contains(t, games, acGame)
}

func createGame(t *testing.T, cl *client.Client, name, token string) *models.Game {
	game, err := cl.CreateGame(&models.Game{
		Name: name,
	}, token)
	require.NoError(t, err)
	require.NotEmpty(t, game.ID)
	require.Empty(t, game.Franchise)
	require.Empty(t, game.FranshiseID)
	require.Equal(t, name, game.Name)

	return game
}
