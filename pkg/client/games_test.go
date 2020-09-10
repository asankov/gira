package client_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/asankov/gira/internal/fixtures"
	"github.com/asankov/gira/pkg/client"
	"github.com/asankov/gira/pkg/models"
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
	if err != nil {
		t.Fatalf("Got non-nil error while constructing client: %v", err)
	}
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
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/games" {
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				if r.Header.Get(models.XAuthToken) != token {
					w.WriteHeader(http.StatusBadRequest)
					return
				}

				w.WriteHeader(http.StatusOK)
				if _, err := w.Write(fixtures.MarshalBytes(t, gameResponse)); err != nil {
					t.Fatalf("error while writing response - %v", err)
				}
			}))
			defer ts.Close()

			cl := newClient(t, ts.URL)

			games, err := cl.GetGames(token, testCase.options)
			if err != nil {
				t.Fatalf("Got non-nil error when calling GetGames: %v", err)
			}

			if len(games) != 1 {
				t.Fatalf("Got %d for length of returned games, expected 1", len(games))
			}

			if games[0].ID != game.ID || games[0].Name != game.Name {
				t.Errorf("Got (%v) for returned game, expected (%v)", games[0], game)
			}
		})
	}
}

func TestGetGamesExcludeAssigned(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/games" || r.Method != http.MethodGet {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if !strings.Contains(r.URL.RawQuery, "excludeAssigned=true") {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if r.Header.Get(models.XAuthToken) != token {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
		if _, err := w.Write(fixtures.MarshalBytes(t, gameResponse)); err != nil {
			t.Fatalf("error while writing response: %v", err)
		}
	}))
	defer ts.Close()

	cl := newClient(t, ts.URL)

	games, err := cl.GetGames(token, &client.GetGamesOptions{ExcludeAssigned: true})
	if err != nil {
		t.Fatalf("Got non-nil error when calling GetGames: %v", err)
	}

	if len(games) != 1 {
		t.Fatalf("Got %d for length of returned games, expected 1", len(games))
	}

	if games[0].ID != game.ID || games[0].Name != game.Name {
		t.Errorf("Got (%v) for returned game, expected (%v)", games[0], game)
	}
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
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(testCase.returnCode)
			}))
			defer ts.Close()

			cl, err := client.New(ts.URL)
			if err != nil {
				t.Fatalf("Got non-nil error while constructing client: %v", err)
			}

			if _, err = cl.GetGames(token, nil); err != testCase.expectedErr {
				t.Fatalf("Got %v error when calling GetGames, expected %v", err, testCase.expectedErr)
			}
		})
	}
}

func TestCreateGame(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/games" || r.Method != http.MethodPost {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if r.Header.Get(models.XAuthToken) != token {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if _, err := w.Write(fixtures.MarshalBytes(t, game)); err != nil {
			t.Fatalf("error while writing response: %v", err)
		}
	}))
	defer ts.Close()

	cl := newClient(t, ts.URL)

	createdGame, err := cl.CreateGame(game, token)
	if err != nil {
		t.Fatalf("Got non-nil error when calling CreateGame: %v", err)
	}

	if createdGame.ID != game.ID || createdGame.Name != game.Name {
		t.Errorf("Got (%v) for created game, expected (%v)", createdGame, game)
	}
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
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(testCase.returnCode)
			}))
			defer ts.Close()

			cl, err := client.New(ts.URL)
			if err != nil {
				t.Fatalf("Got non-nil error while constructing client: %v", err)
			}

			if _, err = cl.CreateGame(game, token); err != testCase.expectedErr {
				t.Fatalf("Got %v error when calling CreateGame, expected %v", err, testCase.expectedErr)
			}
		})
	}
}
