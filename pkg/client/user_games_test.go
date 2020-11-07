package client_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gira-games/api/internal/fixtures"
	"github.com/gira-games/api/pkg/client"
	"github.com/gira-games/api/pkg/models"
	"github.com/google/go-cmp/cmp"
)

var (
	usersGameResponse = map[models.Status][]*models.UserGame{
		models.StatusTODO: {
			{
				ID: "1",
				User: &models.User{
					ID: "1",
				},
				Game: &models.Game{
					ID: "2",
				},
				Status: models.StatusTODO,
			},
		},
		models.StatusInProgress: {
			{
				ID: "2",
				User: &models.User{
					ID: "1",
				},
				Game: &models.Game{
					ID: "3",
				},
				Status: models.StatusInProgress,
			},
		},
	}
)

func TestGetUserGames(t *testing.T) {
	ts := fixtures.NewTestServer(t).
		Path("/users/games").
		Token(token).
		Data(usersGameResponse).
		Build()
	defer ts.Close()

	cl := newClient(t, ts.URL)
	userGames, err := cl.GetUserGames(token)
	assert.NoError(t, err)
	assert.True(t, cmp.Equal(userGames, usersGameResponse))
}

func TestGetUserGameHTTPError(t *testing.T) {
	testCases := []struct {
		name         string
		responseCode int
		expectedErr  error
	}{
		{
			name:         "Auth error",
			responseCode: http.StatusUnauthorized,
			expectedErr:  client.ErrNoAuthorization,
		},
		{
			name:         "Other error",
			responseCode: http.StatusInternalServerError,
			expectedErr:  client.ErrFetchingGames,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			ts := fixtures.NewTestServer(t).
				Path("/users/games").
				Token(token).
				Data(usersGameResponse).
				Return(testCase.responseCode).
				Build()
			defer ts.Close()

			cl := newClient(t, ts.URL)
			_, err := cl.GetUserGames(token)
			assert.Error(t, err, testCase.expectedErr)
		})
	}
}

func TestLinkGameToUser(t *testing.T) {
	ts := fixtures.NewTestServer(t).
		Path("/users/games").
		Method(http.MethodPost).
		Token(token).
		Build()
	defer ts.Close()

	cl := newClient(t, ts.URL)

	_, err := cl.LinkGameToUser("12", token)
	assert.NoError(t, err)
}

func TestLinkGameToUserHTTPError(t *testing.T) {
	testCases := []struct {
		name         string
		responseCode int
		expectedErr  error
	}{
		{
			name:         "Auth error",
			responseCode: http.StatusUnauthorized,
			expectedErr:  client.ErrNoAuthorization,
		},
		{
			name:         "Other error",
			responseCode: http.StatusInternalServerError,
			expectedErr:  client.ErrLinkingGame,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			ts := fixtures.NewTestServer(t).
				Path("/users/games").
				Method(http.MethodPost).
				Token(token).
				Return(testCase.responseCode).
				Build()
			defer ts.Close()

			cl := newClient(t, ts.URL)

			_, err := cl.LinkGameToUser("12", token)
			assert.Error(t, err, testCase.expectedErr)
		})
	}
}

func TestChangeGameStatus(t *testing.T) {
	ts := fixtures.NewTestServer(t).
		Path(fmt.Sprintf("/users/games/%s", game.ID)).
		Method(http.MethodPatch).
		Token(token).
		Build()
	defer ts.Close()

	cl := newClient(t, ts.URL)

	err := cl.ChangeGameStatus(game.ID, token, models.StatusDone)
	assert.NoError(t, err)
}

func TestChangeGameStatusHTTPError(t *testing.T) {
	testCases := []struct {
		name         string
		responseCode int
		expectedErr  error
	}{
		{
			name:         "Auth error",
			responseCode: http.StatusUnauthorized,
			expectedErr:  client.ErrNoAuthorization,
		},
		{
			name:         "Other error",
			responseCode: http.StatusInternalServerError,
			expectedErr:  client.ErrChangingGame,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			ts := fixtures.NewTestServer(t).
				Path(fmt.Sprintf("/users/games/%s", game.ID)).
				Method(http.MethodPatch).
				Token(token).
				Return(testCase.responseCode).
				Build()
			defer ts.Close()

			cl := newClient(t, ts.URL)

			err := cl.ChangeGameStatus(game.ID, token, models.StatusTODO)
			assert.Error(t, err, testCase.expectedErr)
		})
	}
}

func TestChangeGameProgress(t *testing.T) {
	ts := fixtures.NewTestServer(t).
		Path(fmt.Sprintf("/users/games/%s", game.ID)).
		Method(http.MethodPatch).
		Token(token).
		Build()
	defer ts.Close()

	cl := newClient(t, ts.URL)

	err := cl.ChangeGameProgress(game.ID, token, &models.UserGameProgress{
		Current: 10,
		Final:   100,
	})
	assert.NoError(t, err)
}

func TestDeleteGame(t *testing.T) {
	ts := fixtures.NewTestServer(t).
		Path(fmt.Sprintf("/users/games/%s", game.ID)).
		Method(http.MethodDelete).
		Token(token).
		Build()
	defer ts.Close()

	cl := newClient(t, ts.URL)

	err := cl.DeleteUserGame(game.ID, token)
	assert.NoError(t, err)
}

func TestDeleteGameHTTPError(t *testing.T) {
	testCases := []struct {
		name         string
		responseCode int
		expectedErr  error
	}{
		{
			name:         "Auth error",
			responseCode: http.StatusUnauthorized,
			expectedErr:  client.ErrNoAuthorization,
		},
		{
			name:         "Other error",
			responseCode: http.StatusInternalServerError,
			expectedErr:  client.ErrDeletingGame,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			ts := fixtures.NewTestServer(t).
				Path(fmt.Sprintf("/users/games/%s", game.ID)).
				Method(http.MethodDelete).
				Token(token).
				Return(testCase.responseCode).
				Build()
			defer ts.Close()

			cl := newClient(t, ts.URL)

			err := cl.DeleteUserGame(game.ID, token)
			assert.Error(t, err, testCase.expectedErr)
		})
	}
}
