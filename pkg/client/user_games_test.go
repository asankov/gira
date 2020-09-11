package client_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/asankov/gira/internal/fixtures"
	"github.com/asankov/gira/internal/fixtures/assert"
	"github.com/asankov/gira/pkg/client"
	"github.com/asankov/gira/pkg/models"
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
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/users/games" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if r.Header.Get(models.XAuthToken) != token {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
		if _, err := w.Write(fixtures.MarshalBytes(t, usersGameResponse)); err != nil {
			t.Fatalf("error while writing response - %v", err)
		}
	}))
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
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(testCase.responseCode)
			}))
			defer ts.Close()

			cl := newClient(t, ts.URL)
			_, err := cl.GetUserGames(token)
			assert.Error(t, err, testCase.expectedErr)
		})
	}
}
