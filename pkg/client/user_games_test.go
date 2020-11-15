package client_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/asankov/gira/internal/fixtures"
	"github.com/asankov/gira/pkg/client"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
)

var (
	usersGameResponse = map[client.Status][]*client.UserGame{
		"TODO": {
			{
				ID: "1",
				User: &client.User{
					ID: "1",
				},
				Game: &client.Game{
					ID: "2",
				},
				Status: "TODO",
			},
		},
		"In progress": {
			{
				ID: "2",
				User: &client.User{
					ID: "1",
				},
				Game: &client.Game{
					ID: "3",
				},
				Status: "In progress",
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
	resp, err := cl.GetUserGames(context.Background(), &client.GetUserGamesRequest{Token: token})
	assert.NoError(t, err)
	assert.True(t, cmp.Equal(resp.UserGames, usersGameResponse))
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
			_, err := cl.GetUserGames(context.Background(), &client.GetUserGamesRequest{Token: token})
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

	err := cl.LinkGameToUser(context.Background(), &client.LinkGameToUserRequest{
		Token:  token,
		GameID: "12",
	})
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

			err := cl.LinkGameToUser(context.Background(), &client.LinkGameToUserRequest{
				Token:  token,
				GameID: "12",
			})
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

	err := cl.UpdateGameProgress(context.Background(), &client.UpdateGameProgressRequest{
		GameID: game.ID,
		Token:  token,
		Update: client.UpdateGameProgressChange{
			Status: "DONE",
		},
	})
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

			err := cl.UpdateGameProgress(context.Background(), &client.UpdateGameProgressRequest{
				GameID: game.ID,
				Token:  token,
				Update: client.UpdateGameProgressChange{
					Status: "TODO",
				},
			})
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

	err := cl.UpdateGameProgress(context.Background(), &client.UpdateGameProgressRequest{
		GameID: game.ID,
		Token:  token,
		Update: client.UpdateGameProgressChange{
			Status: "DONE",
			Progress: &client.UserGameProgress{
				Current: 10,
				Final:   100,
			},
		},
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

	err := cl.DeleteUserGame(context.Background(), &client.DeleteUserGameRequest{Token: token, GameID: game.ID})
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

			err := cl.DeleteUserGame(context.Background(), &client.DeleteUserGameRequest{Token: token, GameID: game.ID})
			assert.Error(t, err, testCase.expectedErr)
		})
	}
}
