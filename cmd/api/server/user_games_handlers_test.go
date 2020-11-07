package server

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	gassert "github.com/gira-games/api/internal/fixtures/assert"

	"github.com/gira-games/api/internal/fixtures"
	"github.com/gira-games/api/pkg/models"
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

	gassert.StatusOK(t, w)

	var gamesResponse map[models.Status][]*models.UserGame
	fixtures.Decode(t, w.Body, &gamesResponse)

	require.Equal(t, len(gamesResponse), len(expectedGames))
	for _, g := range gamesResponse["To Do"] {
		assert.True(t, gameIn(g, expectedGames["To Do"]))
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

	gassert.StatusCode(t, w, http.StatusInternalServerError)
}

func TestUserGamesPost(t *testing.T) {
	testCases := []struct {
		Name              string
		ProgressInRequest *models.UserGameProgress
		ExpectedProgress  *models.UserGameProgress
	}{
		{
			Name: "Progress.Final is 100 and Progress.Current is 0, if nothing is passed",
			ExpectedProgress: &models.UserGameProgress{
				Current: 0,
				Final:   100,
			},
		},
		{
			Name: "Progress.Final and Progress.Current are equal to what is passed",
			ProgressInRequest: &models.UserGameProgress{
				Current: 2,
				Final:   50,
			},
			ExpectedProgress: &models.UserGameProgress{
				Current: 2,
				Final:   50,
			},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
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
				LinkGameToUser(gomock.Eq("12"), gomock.Eq(gameID), gomock.Eq(testCase.ExpectedProgress)).
				Return(nil)

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, "/users/games", fixtures.Marshal(t, &models.UserGameRequest{Game: &models.Game{ID: gameID}, Progress: testCase.ProgressInRequest}))
			r.Header.Add(models.XAuthToken, token)

			srv.ServeHTTP(w, r)

			gassert.StatusOK(t, w)
		})
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
		LinkGameToUser(gomock.Eq("12"), gomock.Eq(gameID), gomock.Any()).
		Return(errors.New("error returned on purpose"))

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/users/games", fixtures.Marshal(t, &models.UserGameRequest{Game: &models.Game{ID: gameID}}))
	r.Header.Add(models.XAuthToken, token)

	srv.ServeHTTP(w, r)

	gassert.StatusCode(t, w, http.StatusInternalServerError)
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

	gassert.StatusCode(t, w, http.StatusBadRequest)
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
	r := httptest.NewRequest(http.MethodPatch, "/users/games/1", fixtures.Marshal(t, models.ChangeGameStatusRequest{Status: models.StatusDone}))
	r.Header.Add(models.XAuthToken, token)

	srv.ServeHTTP(w, r)

	gassert.StatusOK(t, w)
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

	gassert.StatusCode(t, w, http.StatusBadRequest)
}

func TestUsersGamesPatchServiceError(t *testing.T) {
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
	r := httptest.NewRequest(http.MethodPatch, "/users/games/1", fixtures.Marshal(t, models.ChangeGameStatusRequest{Status: models.StatusDone}))
	r.Header.Add(models.XAuthToken, token)

	srv.ServeHTTP(w, r)

	gassert.StatusCode(t, w, http.StatusInternalServerError)
}

func TestUsersGamesDelete(t *testing.T) {
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
	userGamesModelMock.EXPECT().GetUserGames(gomock.Eq("12")).Return([]*models.UserGame{
		{
			ID: "1",
		},
	}, nil)
	userGamesModelMock.
		EXPECT().
		DeleteUserGame(gomock.Eq("1")).
		Return(nil)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodDelete, "/users/games/1", nil)
	r.Header.Add(models.XAuthToken, token)

	srv.ServeHTTP(w, r)

	gassert.StatusOK(t, w)
}

func TestUsersGamesDeleteUserDoesNotOwnGame(t *testing.T) {
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
	userGamesModelMock.EXPECT().GetUserGames(gomock.Eq("12")).Return([]*models.UserGame{
		{
			ID: "123",
		},
	}, nil)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodDelete, "/users/games/1", nil)
	r.Header.Add(models.XAuthToken, token)

	srv.ServeHTTP(w, r)

	gassert.StatusCode(t, w, http.StatusBadRequest)
}

func TestUsersGamesDeleteServiceError(t *testing.T) {
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
	userGamesModelMock.EXPECT().GetUserGames(gomock.Eq("12")).Return([]*models.UserGame{
		{
			ID: "1",
		},
	}, nil)
	userGamesModelMock.
		EXPECT().
		DeleteUserGame(gomock.Eq("1")).
		Return(errors.New("error while deleting game"))

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodDelete, "/users/games/1", nil)
	r.Header.Add(models.XAuthToken, token)

	srv.ServeHTTP(w, r)

	gassert.StatusCode(t, w, http.StatusInternalServerError)
}
