package server

import (
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/asankov/gira/internal/fixtures"
	"github.com/asankov/gira/pkg/models"
	"github.com/golang/mock/gomock"
)

func setupUserGamesServer(u UserGamesModel, a *fixtures.AuthenticatorMock) *Server {
	return &Server{
		Log:            log.New(os.Stdout, "", 0),
		UserGamesModel: u,
		Authenticator:  a,
	}
}

func TestUsersGamesGet(t *testing.T) {
	ctrl := gomock.NewController(t)

	authenticatorMock := fixtures.NewAuthenticatorMock(ctrl)
	userGamesMock := fixtures.NewUserGamesModelMock(ctrl)

	authenticatorMock.EXPECT().
		DecodeToken(gomock.Eq(token)).
		Return(&models.User{
			ID: "12",
		}, nil)

	expectedGames := []*models.Game{
		{ID: "1", Name: "ACI"},
		{ID: "2", Name: "ACII"},
	}
	userGamesMock.EXPECT().
		GetUserGames(gomock.Eq("12")).
		Return(expectedGames, nil)

	srv := setupUserGamesServer(userGamesMock, authenticatorMock)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/users/12/games", nil)
	r.Header.Add("x-auth-token", token)

	srv.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("Got (%d) for HTTP StatusCode, expected (%d)", w.Code, http.StatusOK)
	}

	var games []*models.Game
	fixtures.Decode(t, w.Body, &games)

	got, expected := len(games), len(expectedGames)
	if got != expected {
		t.Errorf("Got (%d) for length of returned games, expected (%d)", got, expected)
	}
	for _, g := range games {
		if !gameIn(g, expectedGames) {
			t.Errorf("Expected game (%#v) to be in returned games", g)
		}
	}
}

func gameIn(game *models.Game, games []*models.Game) bool {
	for _, g := range games {
		if game.ID == g.ID && game.Name == g.Name {
			return true
		}
	}
	return false
}
