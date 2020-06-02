package server

import (
	"bytes"
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/asankov/gira/internal/fixtures"
	"github.com/asankov/gira/pkg/models"
	"github.com/golang/mock/gomock"
)

// TODO: this is the same in all test files.
// remove in favour of server.New(*ServerOptions) method
func setupUserGamesServer(ug UserGamesModel, u UserModel, a *fixtures.AuthenticatorMock) *Server {
	return &Server{
		Log:            log.New(os.Stdout, "", 0),
		UserGamesModel: ug,
		UserModel:      u,
		Authenticator:  a,
	}
}

func TestUsersGamesGet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	authenticatorMock := fixtures.NewAuthenticatorMock(ctrl)
	userGamesModelMock := fixtures.NewUserGamesModelMock(ctrl)
	userModelMock := fixtures.NewUserModelMock(ctrl)
	srv := setupUserGamesServer(userGamesModelMock, userModelMock, authenticatorMock)

	authenticatorMock.EXPECT().
		DecodeToken(gomock.Eq(token)).
		Return(nil, nil)
	userModelMock.EXPECT().
		GetUserByToken(gomock.Eq(token)).
		Return(&models.User{
			ID: "12",
		}, nil)

	expectedGames := map[models.Status][]*models.UserGame{
		"To Do": {
			{ID: "1", Game: &models.Game{ID: "1", Name: "ACI"}},
			{ID: "2", Game: &models.Game{ID: "2", Name: "ACII"}},
		},
	}

	userGamesModelMock.EXPECT().
		GetUserGamesGrouped(gomock.Eq("12")).
		Return(expectedGames, nil)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/users/games", nil)
	r.Header.Add("x-auth-token", token)

	srv.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("Got (%d) for HTTP StatusCode, expected (%d)", w.Code, http.StatusOK)
	}

	var gamesResponse map[models.Status][]*models.UserGame
	fixtures.Decode(t, w.Body, &gamesResponse)

	got, expected := len(gamesResponse), len(expectedGames)
	if got != expected {
		t.Errorf("Got (%d) for length of returned games, expected (%d)", got, expected)
	}
	for _, g := range gamesResponse["To Do"] {
		if !gameIn(g, expectedGames["To Do"]) {
			t.Errorf("Expected game (%#v) to be in returned games", g)
		}
	}
}

func gameIn(game *models.UserGame, games []*models.UserGame) bool {
	for _, g := range games {
		if game.ID == g.ID && game.Game.Name == g.Game.Name && game.Game.ID == g.Game.ID {
			return true
		}
	}
	return false
}

func TestUsersGamesGetInternalError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	authenticatorMock := fixtures.NewAuthenticatorMock(ctrl)
	userGamesModelMock := fixtures.NewUserGamesModelMock(ctrl)
	userModelMock := fixtures.NewUserModelMock(ctrl)
	srv := setupUserGamesServer(userGamesModelMock, userModelMock, authenticatorMock)

	authenticatorMock.EXPECT().
		DecodeToken(gomock.Eq(token)).
		Return(nil, nil)
	userModelMock.EXPECT().
		GetUserByToken(gomock.Eq(token)).
		Return(&models.User{
			ID: "12",
		}, nil)

	userGamesModelMock.EXPECT().
		GetUserGamesGrouped(gomock.Eq("12")).
		Return(nil, errors.New("error returned on purpose"))

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/users/games", nil)
	r.Header.Add("x-auth-token", token)

	srv.ServeHTTP(w, r)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("Got (%d) for HTTP StatusCode, expected (%d)", w.Code, http.StatusInternalServerError)
	}
}

func TestUserGamesPost(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	authenticatorMock := fixtures.NewAuthenticatorMock(ctrl)
	userGamesModelMock := fixtures.NewUserGamesModelMock(ctrl)
	userModelMock := fixtures.NewUserModelMock(ctrl)
	srv := setupUserGamesServer(userGamesModelMock, userModelMock, authenticatorMock)

	authenticatorMock.EXPECT().
		DecodeToken(gomock.Eq(token)).
		Return(nil, nil)
	userModelMock.EXPECT().
		GetUserByToken(gomock.Eq(token)).
		Return(&models.User{
			ID: "12",
		}, nil)

	gameID := "666"
	userGamesModelMock.EXPECT().
		LinkGameToUser(gomock.Eq("12"), gomock.Eq(gameID)).
		Return(nil)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/users/games", fixtures.Marshall(t, &models.UserGameRequest{Game: &models.Game{ID: gameID}}))
	r.Header.Add("x-auth-token", token)

	srv.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("Got (%d) for HTTP StatusCode, expected (%d)", w.Code, http.StatusOK)
	}
}

func TestUsersGamesPostInternalError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	authenticatorMock := fixtures.NewAuthenticatorMock(ctrl)
	userGamesModelMock := fixtures.NewUserGamesModelMock(ctrl)
	userModelMock := fixtures.NewUserModelMock(ctrl)
	srv := setupUserGamesServer(userGamesModelMock, userModelMock, authenticatorMock)

	authenticatorMock.EXPECT().
		DecodeToken(gomock.Eq(token)).
		Return(nil, nil)
	userModelMock.EXPECT().
		GetUserByToken(gomock.Eq(token)).
		Return(&models.User{
			ID: "12",
		}, nil)

	gameID := "666"
	userGamesModelMock.EXPECT().
		LinkGameToUser(gomock.Eq("12"), gomock.Eq(gameID)).
		Return(errors.New("error returned on purpose"))

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/users/games", fixtures.Marshall(t, &models.UserGameRequest{Game: &models.Game{ID: gameID}}))
	r.Header.Add("x-auth-token", token)

	srv.ServeHTTP(w, r)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("Got (%d) for HTTP StatusCode, expected (%d)", w.Code, http.StatusInternalServerError)
	}
}

func TestUsersGamesPostParseError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	authenticatorMock := fixtures.NewAuthenticatorMock(ctrl)
	userGamesModelMock := fixtures.NewUserGamesModelMock(ctrl)
	userModelMock := fixtures.NewUserModelMock(ctrl)
	srv := setupUserGamesServer(userGamesModelMock, userModelMock, authenticatorMock)

	authenticatorMock.EXPECT().
		DecodeToken(gomock.Eq(token)).
		Return(nil, nil)
	userModelMock.EXPECT().
		GetUserByToken(gomock.Eq(token)).
		Return(&models.User{
			ID: "12",
		}, nil)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/users/games", bytes.NewBuffer(nil))
	r.Header.Add("x-auth-token", token)

	srv.ServeHTTP(w, r)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("Got (%d) for HTTP StatusCode, expected (%d)", w.Code, http.StatusBadRequest)
	}
}
