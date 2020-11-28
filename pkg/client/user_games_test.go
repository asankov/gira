package client_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/asankov/gira/internal/fixtures"
	"github.com/asankov/gira/pkg/client"
	"github.com/stretchr/testify/assert"
)

func TestChangeGameStatus(t *testing.T) {
	ts := fixtures.NewTestServer(t).
		Path(fmt.Sprintf("/games/%s", game.ID)).
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
				Path(fmt.Sprintf("/games/%s", game.ID)).
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
		Path(fmt.Sprintf("/games/%s", game.ID)).
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
			Progress: &client.GameProgress{
				Current: 10,
				Final:   100,
			},
		},
	})
	assert.NoError(t, err)
}

func TestDeleteGame(t *testing.T) {
	ts := fixtures.NewTestServer(t).
		Path(fmt.Sprintf("/games/%s", game.ID)).
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
				Path(fmt.Sprintf("/games/%s", game.ID)).
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
