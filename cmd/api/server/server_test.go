package server

import (
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
)

func TestGetGames(t *testing.T) {
	ctrl := gomock.NewController(t)

	gameModel := fixtures.NewGameModelMock(ctrl)
	srv := Server{
		Log:       log.New(os.Stdout, "", 0),
		GameModel: gameModel,
		Auth:      authenticator,
	}

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
	token, err := srv.Auth.NewTokenForUser("anton")
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
	if err := json.NewDecoder(w.Body).Decode(&res); err != nil {
		t.Fatalf("Got unexpected error while decoding response - %v", err)
	}

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
	srv := Server{
		Log:       log.New(os.Stdout, "", 0),
		GameModel: gameModel,
		Auth:      authenticator,
	}

	gameModel.
		EXPECT().
		All().
		Return(nil, errors.New("this is an intentional error"))

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/games", nil)
	token, err := srv.Auth.NewTokenForUser("anton")
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
