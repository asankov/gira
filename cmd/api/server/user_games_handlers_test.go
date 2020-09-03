package server

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/asankov/gira/internal/fixtures"
	"github.com/asankov/gira/pkg/models"
	"github.com/golang/mock/gomock"
)

func TestUsersGamesGet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	authenticatorMock := fixtures.NewAuthenticatorMock(ctrl)
	userGamesModelMock := fixtures.NewUserGamesModelMock(ctrl)
	userModelMock := fixtures.NewUserModelMock(ctrl)
	srv := newServer(t, &Options{
		Authenticator:  authenticatorMock,
		UserModel:      userModelMock,
		UserGamesModel: userGamesModelMock,
	})

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
	r.Header.Add(models.XAuthToken, token)

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
	srv := newServer(t, &Options{
		Authenticator:  authenticatorMock,
		UserModel:      userModelMock,
		UserGamesModel: userGamesModelMock,
	})

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
	r.Header.Add(models.XAuthToken, token)

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
	srv := newServer(t, &Options{
		Authenticator:  authenticatorMock,
		UserModel:      userModelMock,
		UserGamesModel: userGamesModelMock,
	})
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
	r.Header.Add(models.XAuthToken, token)

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
	srv := newServer(t, &Options{
		Authenticator:  authenticatorMock,
		UserModel:      userModelMock,
		UserGamesModel: userGamesModelMock,
	})
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
	r.Header.Add(models.XAuthToken, token)

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
	srv := newServer(t, &Options{
		Authenticator:  authenticatorMock,
		UserModel:      userModelMock,
		UserGamesModel: userGamesModelMock,
	})
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
	r.Header.Add(models.XAuthToken, token)

	srv.ServeHTTP(w, r)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("Got (%d) for HTTP StatusCode, expected (%d)", w.Code, http.StatusBadRequest)
	}
}

func TestUsersGamesPatch(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	authenticatorMock := fixtures.NewAuthenticatorMock(ctrl)
	userGamesModelMock := fixtures.NewUserGamesModelMock(ctrl)
	userModelMock := fixtures.NewUserModelMock(ctrl)
	srv := newServer(t, &Options{
		Authenticator:  authenticatorMock,
		UserModel:      userModelMock,
		UserGamesModel: userGamesModelMock,
	})

	authenticatorMock.EXPECT().
		DecodeToken(gomock.Eq(token)).
		Return(nil, nil)
	userModelMock.EXPECT().
		GetUserByToken(gomock.Eq(token)).
		Return(&models.User{
			ID: "12",
		}, nil)
	userGamesModelMock.
		EXPECT().
		ChangeGameStatus(gomock.Eq("12"), gomock.Eq("1"), gomock.Eq(models.StatusDone)).
		Return(nil)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPatch, "/users/games/1", fixtures.Marshall(t, models.ChangeGameStatusRequest{Status: models.StatusDone}))
	r.Header.Add(models.XAuthToken, token)

	srv.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("Got (%d) for HTTP StatusCode, expected (%d)", w.Code, http.StatusOK)
	}
}

func TestUsersGamesPatchInvalidBody(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	authenticatorMock := fixtures.NewAuthenticatorMock(ctrl)
	userModelMock := fixtures.NewUserModelMock(ctrl)
	srv := newServer(t, &Options{
		Authenticator: authenticatorMock,
		UserModel:     userModelMock,
	})

	authenticatorMock.EXPECT().
		DecodeToken(gomock.Eq(token)).
		Return(nil, nil)
	userModelMock.EXPECT().
		GetUserByToken(gomock.Eq(token)).
		Return(&models.User{
			ID: "12",
		}, nil)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPatch, "/users/games/1", nil)
	r.Header.Add(models.XAuthToken, token)

	srv.ServeHTTP(w, r)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("Got (%d) for HTTP StatusCode, expected (%d)", w.Code, http.StatusBadRequest)
	}
}

func TestUsersGamesPatchErrChangingStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	authenticatorMock := fixtures.NewAuthenticatorMock(ctrl)
	userGamesModelMock := fixtures.NewUserGamesModelMock(ctrl)
	userModelMock := fixtures.NewUserModelMock(ctrl)
	srv := newServer(t, &Options{
		Authenticator:  authenticatorMock,
		UserModel:      userModelMock,
		UserGamesModel: userGamesModelMock,
	})

	authenticatorMock.EXPECT().
		DecodeToken(gomock.Eq(token)).
		Return(nil, nil)
	userModelMock.EXPECT().
		GetUserByToken(gomock.Eq(token)).
		Return(&models.User{
			ID: "12",
		}, nil)
	userGamesModelMock.
		EXPECT().
		ChangeGameStatus(gomock.Eq("12"), gomock.Eq("1"), gomock.Eq(models.StatusDone)).
		Return(errors.New("error while changing game status"))

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPatch, "/users/games/1", fixtures.Marshall(t, models.ChangeGameStatusRequest{Status: models.StatusDone}))
	r.Header.Add(models.XAuthToken, token)

	srv.ServeHTTP(w, r)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("Got (%d) for HTTP StatusCode, expected (%d)", w.Code, http.StatusInternalServerError)
	}
}
