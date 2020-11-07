package client_test

import (
	"net/http"
	"testing"

	gassert "github.com/gira-games/api/internal/fixtures/assert"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/gira-games/api/internal/fixtures"
	"github.com/gira-games/api/pkg/client"
	"github.com/gira-games/api/pkg/models"
)

var (
	token = "my-token"
	game  = &models.Game{
		ID:   "1",
		Name: "A",
	}
	gameResponse = models.GamesResponse{
		Games: []*models.Game{game},
	}
)

func newClient(t *testing.T, url string) *client.Client {
	cl, err := client.New(url)
	require.NoError(t, err)

	return cl
}

func TestGetGames(t *testing.T) {
	testCases := []struct {
		name    string
		options *client.GetGamesOptions
	}{
		{
			name:    "Empty options",
			options: &client.GetGamesOptions{},
		},
		{
			name:    "Nil options",
			options: nil,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			ts := fixtures.NewTestServer(t).
				Path("/games").
				Data(gameResponse).
				Token(token).
				Build()
			defer ts.Close()

			cl := newClient(t, ts.URL)

			games, err := cl.GetGames(token, testCase.options)

			require.NoError(t, err)
			assert.Equal(t, 1, len(games))
			assert.Equal(t, game.ID, games[0].ID)
			assert.Equal(t, game.Name, games[0].Name)
		})
	}
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

	games, err := cl.GetGames(token, &client.GetGamesOptions{ExcludeAssigned: true})

	require.NoError(t, err)
	assert.Equal(t, 1, len(games))
	assert.Equal(t, game.ID, games[0].ID)
	assert.Equal(t, game.Name, games[0].Name)
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

			games, err := cl.GetGames(token, nil)
			assert.Nil(t, games)
			gassert.Error(t, err, testCase.expectedErr)
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

	createdGame, err := cl.CreateGame(game, token)

	require.NoError(t, err)
	assert.Equal(t, game.ID, createdGame.ID)
	assert.Equal(t, game.Name, createdGame.Name)
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

			createdGame, err := cl.CreateGame(game, token)
			assert.Nil(t, createdGame)
			gassert.Error(t, err, testCase.expectedErr)
		})
	}
}
