package server

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/asankov/gira/internal/auth"
	"github.com/asankov/gira/internal/fixtures"
	"github.com/asankov/gira/pkg/models"
	"github.com/golang/mock/gomock"
)

var (
	authenticator = auth.NewAutheniticator("test_secret")
	user          = &models.User{Username: "anton"}
)

func newServer(g GameModel, a *auth.Authenticator) *Server {
	return &Server{
		Log:       log.New(os.Stdout, "", 0),
		GameModel: g,
		Auth:      a,
	}
}
func TestGetGames(t *testing.T) {
	ctrl := gomock.NewController(t)

	gameModel := fixtures.NewGameModelMock(ctrl)
	srv := newServer(gameModel, authenticator)

	gamesResponse := []*models.Game{
		{ID: "1", Name: "AC"},
		{ID: "2", Name: "ACII"},
	}
	gameModel.
		EXPECT().
		All().
		Return(gamesResponse, nil)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/games", nil)
	token, err := srv.Auth.NewTokenForUser(user)
	if err != nil {
		t.Fatalf("Got unexpected error while trying to generate token for user - %v", err)
	}
	r.Header.Set("x-auth-token", token)
	srv.ServeHTTP(w, r)

	got, expected := w.Result().StatusCode, http.StatusOK
	if got != expected {
		t.Fatalf("Got status code - (%d), expected (%d)", got, expected)
	}

	var res []*models.Game
	decode(t, w, &res)

	if len(res) != 2 {
		t.Fatalf("Got (%d) for length of result, expected %d", len(res), 2)
	}
	for i := 0; i < len(res); i++ {
		got, expected := res[i].ID, gamesResponse[i].ID
		if got != expected {
			t.Fatalf("Got (%s) for result[%d].ID, expected (%s)", got, i, expected)
		}
		got, expected = res[i].Name, gamesResponse[i].Name
		if got != expected {
			t.Fatalf("Got (%s) for result[%d].Name, expected (%s)", got, i, expected)
		}
	}
}

func TestGetGamesErr(t *testing.T) {
	ctrl := gomock.NewController(t)

	gameModel := fixtures.NewGameModelMock(ctrl)
	srv := newServer(gameModel, authenticator)

	gameModel.
		EXPECT().
		All().
		Return(nil, errors.New("this is an intentional error"))

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/games", nil)
	token, err := srv.Auth.NewTokenForUser(user)
	if err != nil {
		t.Fatalf("Got unexpected error while trying to generate token for user - %v", err)
	}
	r.Header.Set("x-auth-token", token)
	srv.ServeHTTP(w, r)

	got, expected := w.Result().StatusCode, http.StatusInternalServerError
	if got != expected {
		t.Fatalf("Got status code - (%d), expected (%d)", got, expected)
	}
}

func TestCreateGame(t *testing.T) {
	ctrl := gomock.NewController(t)

	gameModel := fixtures.NewGameModelMock(ctrl)
	srv := newServer(gameModel, authenticator)

	actualName := "ACIII"
	actualGame := &models.Game{Name: actualName}
	gameModel.
		EXPECT().
		Insert(actualGame).
		Return(actualGame, nil)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/games", marshall(t, actualGame))
	token, err := srv.Auth.NewTokenForUser(user)
	if err != nil {
		t.Fatalf("Got unexpected error while trying to generate token for user - %v", err)
	}
	r.Header.Set("x-auth-token", token)
	srv.ServeHTTP(w, r)

	statusCode := w.Result().StatusCode
	if statusCode != http.StatusOK {
		t.Fatalf("Got (%d) for status code, expected (%d)", statusCode, http.StatusOK)
	}

	var game *models.Game
	decode(t, w, &game)
	if game.Name != actualName {
		t.Fatalf("Got (%s) for game.Name, expected (%s)", game.Name, actualName)
	}
}

func TestCreateGameError(t *testing.T) {

	cases := []struct {
		name string
		game *models.Game
	}{
		{
			name: "Empty game",
			game: &models.Game{Name: ""},
		},
		{
			name: "Filled ID",
			game: &models.Game{ID: "123", Name: "something valid"},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			gameModel := fixtures.NewGameModelMock(ctrl)
			srv := newServer(gameModel, authenticator)

			actualGame := c.game
			gameModel.
				EXPECT().
				Insert(actualGame).
				Return(actualGame, nil)

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, "/games", marshall(t, actualGame))
			token, err := srv.Auth.NewTokenForUser(user)
			if err != nil {
				t.Fatalf("Got unexpected error while trying to generate token for user - %v", err)
			}
			r.Header.Set("x-auth-token", token)
			srv.ServeHTTP(w, r)

			got, expected := w.Result().StatusCode, http.StatusBadRequest
			if got != expected {
				t.Fatalf("Got (%d) for status code, expected (%d)", got, expected)
			}
		})
	}
}

func marshall(t *testing.T, payload interface{}) *bytes.Buffer {
	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("Got unexpected error while marshalling payload - %v", err)
	}
	return bytes.NewBuffer(body)
}

func decode(t *testing.T, w *httptest.ResponseRecorder, into interface{}) {
	if err := json.NewDecoder(w.Body).Decode(&into); err != nil {
		t.Fatalf("Got unexpected error while decoding response - %v", err)
	}
}
