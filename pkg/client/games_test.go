package client_test

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/asankov/gira/internal/fixtures"
	"github.com/asankov/gira/pkg/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	token = "my-token"
	game  = &client.Game{
		ID:   "1",
		Name: "A",
	}
	gameResponse = client.GetGamesResponse{
		Games: []*client.Game{game},
	}
)

func newClient(t *testing.T, url string) *client.Client {
	cl, err := client.New(url)
	require.NoError(t, err)

	return cl
}

func TestGetGames(t *testing.T) {
	ts := fixtures.NewTestServer(t).
		Path("/games").
		Data(gameResponse).
		Token(token).
		Build()
	defer ts.Close()

	cl := newClient(t, ts.URL)

	games, err := cl.GetGames(context.Background(), &client.GetGamesRequest{Token: token})

	require.NoError(t, err)
	require.Equal(t, 1, len(games.Games))
	assert.Equal(t, game.ID, games.Games[0].ID)
	assert.Equal(t, game.Name, games.Games[0].Name)
}

func TestGetGamesExcludeAssigned(t *testing.T) {
	ts := fixtures.NewTestServer(t).
		Path("/games").
		Data(gameResponse).
		Token(token).
		Query("excludeAssigned=true").
		Build()
	defer ts.Close()

	cl := newClient(t, ts.URL)

	games, err := cl.GetGames(context.Background(), &client.GetGamesRequest{Token: token, ExcludeAssigned: true})

	require.NoError(t, err)
	require.Equal(t, 1, len(games.Games))
	assert.Equal(t, game.ID, games.Games[0].ID)
	assert.Equal(t, game.Name, games.Games[0].Name)
}

func TestGetGamesHTTPError(t *testing.T) {
	testCases := []struct {
		name        string
		returnCode  int
		expectedErr error
	}{
		{
			name:        "Auth error",
			returnCode:  http.StatusUnauthorized,
			expectedErr: client.ErrNoAuthorization,
		},
		{
			name:        "Other error",
			returnCode:  http.StatusBadRequest,
			expectedErr: client.ErrFetchingGames,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			ts := fixtures.NewTestServer(t).
				Path("/games").
				Return(testCase.returnCode).
				Build()
			defer ts.Close()

			cl, err := client.New(ts.URL)
			require.NoError(t, err)

			games, err := cl.GetGames(context.Background(), &client.GetGamesRequest{Token: token})
			assert.Nil(t, games)
			assert.True(t, errors.Is(err, testCase.expectedErr))
		})
	}
}

func TestCreateGame(t *testing.T) {
	ts := fixtures.NewTestServer(t).
		Path("/games").
		Method(http.MethodPost).
		Data(game).
		Return(http.StatusOK).
		Token(token).
		Build()
	defer ts.Close()

	cl := newClient(t, ts.URL)

	res, err := cl.CreateGame(context.Background(), &client.CreateGameRequest{Token: token, Game: game})

	require.NoError(t, err)
	assert.Equal(t, game.ID, res.Game.ID)
	assert.Equal(t, game.Name, res.Game.Name)
	assert.Empty(t, res.Game.FranchiseID)
}

func TestCreateGameHTTPError(t *testing.T) {
	testCases := []struct {
		name        string
		returnCode  int
		expectedErr error
	}{
		{
			name:        "Auth error",
			returnCode:  http.StatusUnauthorized,
			expectedErr: client.ErrNoAuthorization,
		},
		{
			name:        "Other error",
			returnCode:  http.StatusBadRequest,
			expectedErr: client.ErrCreatingGame,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			ts := fixtures.NewTestServer(t).
				Path("/games").
				Method(http.MethodPost).
				Return(testCase.returnCode).
				Build()
			defer ts.Close()

			cl, err := client.New(ts.URL)
			require.NoError(t, err)

			createdGame, err := cl.CreateGame(context.Background(), &client.CreateGameRequest{Token: token, Game: game})
			assert.Nil(t, createdGame)
			assert.True(t, errors.Is(err, testCase.expectedErr))
		})
	}
}
