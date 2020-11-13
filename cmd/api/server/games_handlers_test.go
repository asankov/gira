package server

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	gassert "github.com/gira-games/api/internal/fixtures/assert"

	"github.com/gira-games/api/internal/fixtures"
	"github.com/gira-games/api/pkg/models"
	"github.com/gira-games/api/pkg/models/postgres"
	"github.com/golang/mock/gomock"
)

var (
	user = &models.User{Username: "anton"}
)

func TestGetGames(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	gameModel := fixtures.NewGameModelMock(ctrl)
	userModel := fixtures.NewUserModelMock(ctrl)
	authenticator := fixtures.NewAuthenticatorMock(ctrl)
	srv := newServer(t, &Options{
		Authenticator: authenticator,
		GameModel:     gameModel,
		UserModel:     userModel,
	})

	gamesResponse := []*models.Game{
		{ID: "1", Name: "AC"},
		{ID: "2", Name: "ACII"},
	}
	authenticator.EXPECT().
		DecodeToken(gomock.Eq(token)).
		Return(user, nil)
	gameModel.
		EXPECT().
		AllForUser(user.ID).
		Return(gamesResponse, nil)
	userModel.
		EXPECT().
		GetUserByToken(token).
		Return(user, nil)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/games", nil)
	r.Header.Set(models.XAuthToken, token)
	srv.ServeHTTP(w, r)

	gassert.StatusOK(t, w)

	var res models.GamesResponse
	fixtures.Decode(t, w.Body, &res)

	require.Equal(t, 2, len(res.Games))
	for i := 0; i < len(res.Games); i++ {
		assert.Equal(t, gamesResponse[i].ID, res.Games[i].ID)
		assert.Equal(t, gamesResponse[i].Name, res.Games[i].Name)
	}
}

func TestGetGamesErr(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	gameModel := fixtures.NewGameModelMock(ctrl)
	userModel := fixtures.NewUserModelMock(ctrl)
	authenticator := fixtures.NewAuthenticatorMock(ctrl)
	srv := newServer(t, &Options{
		Authenticator: authenticator,
		GameModel:     gameModel,
		UserModel:     userModel,
	})
	authenticator.EXPECT().
		DecodeToken(gomock.Eq(token)).
		Return(user, nil)
	gameModel.
		EXPECT().
		AllForUser(user.ID).
		Return(nil, errors.New("this is an intentional error"))
	userModel.
		EXPECT().
		GetUserByToken(token).
		Return(user, nil)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/games", nil)
	r.Header.Set(models.XAuthToken, token)
	srv.ServeHTTP(w, r)

	gassert.StatusCode(t, w, http.StatusInternalServerError)
}

func TestGetGameByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	gameModel := fixtures.NewGameModelMock(ctrl)
	userModel := fixtures.NewUserModelMock(ctrl)
	authenticator := fixtures.NewAuthenticatorMock(ctrl)
	srv := newServer(t, &Options{
		Authenticator: authenticator,
		GameModel:     gameModel,
		UserModel:     userModel,
	})
	actualName := "ACIII"
	actualGame := &models.Game{Name: actualName}
	authenticator.EXPECT().
		DecodeToken(gomock.Eq(token)).
		Return(user, nil)
	gameModel.
		EXPECT().
		Get("1").
		Return(actualGame, nil)
	userModel.
		EXPECT().
		GetUserByToken(token).
		Return(user, nil)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/games/1", nil)
	r.Header.Set(models.XAuthToken, token)
	srv.ServeHTTP(w, r)

	gassert.StatusOK(t, w)

	var game *models.Game
	fixtures.Decode(t, w.Body, &game)

	assert.Equal(t, game.Name, actualGame.Name)
}

func TestGetGameByIDDBError(t *testing.T) {
	cases := []struct {
		name         string
		dbError      error
		expectedCode int
	}{
		{
			name:         "Error no record",
			dbError:      postgres.ErrNoRecord,
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "Other error",
			dbError:      errors.New("some unknown error"),
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			gameModel := fixtures.NewGameModelMock(ctrl)
			userModel := fixtures.NewUserModelMock(ctrl)
			authenticator := fixtures.NewAuthenticatorMock(ctrl)
			srv := newServer(t, &Options{
				Authenticator: authenticator,
				GameModel:     gameModel,
				UserModel:     userModel,
			})
			authenticator.EXPECT().
				DecodeToken(gomock.Eq(token)).
				Return(user, nil)
			gameModel.
				EXPECT().
				Get("1").
				Return(nil, c.dbError)
			userModel.
				EXPECT().
				GetUserByToken(token).
				Return(user, nil)

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/games/1", nil)
			r.Header.Set(models.XAuthToken, token)
			srv.ServeHTTP(w, r)

			gassert.StatusCode(t, w, c.expectedCode)
		})
	}

}
func TestCreateGame(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	gameModel := fixtures.NewGameModelMock(ctrl)
	userModel := fixtures.NewUserModelMock(ctrl)
	authenticator := fixtures.NewAuthenticatorMock(ctrl)
	srv := newServer(t, &Options{
		Authenticator: authenticator,
		GameModel:     gameModel,
		UserModel:     userModel,
	})

	authenticator.EXPECT().
		DecodeToken(gomock.Eq(token)).
		Return(user, nil)

	actualName := "ACIII"
	actualGame := &models.Game{Name: actualName}
	gameModel.
		EXPECT().
		Insert(actualGame).
		Return(actualGame, nil)
	userModel.
		EXPECT().
		GetUserByToken(token).
		Return(user, nil)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/games", fixtures.Marshal(t, actualGame))
	r.Header.Set(models.XAuthToken, token)
	srv.ServeHTTP(w, r)

	gassert.StatusOK(t, w)

	var game *models.Game
	fixtures.Decode(t, w.Body, &game)
	assert.Equal(t, actualName, game.Name)
}

func TestCreateGameValidationError(t *testing.T) {

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
			defer ctrl.Finish()

			gameModel := fixtures.NewGameModelMock(ctrl)
			userModel := fixtures.NewUserModelMock(ctrl)
			authenticator := fixtures.NewAuthenticatorMock(ctrl)
			srv := newServer(t, &Options{
				Authenticator: authenticator,
				GameModel:     gameModel,
				UserModel:     userModel,
			})
			authenticator.EXPECT().
				DecodeToken(gomock.Eq(token)).
				Return(user, nil)
			userModel.
				EXPECT().
				GetUserByToken(token).
				Return(user, nil)

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, "/games", fixtures.Marshal(t, c.game))
			r.Header.Set(models.XAuthToken, token)
			srv.ServeHTTP(w, r)

			gassert.StatusCode(t, w, http.StatusBadRequest)
		})
	}
}

func TestCreateGameDBError(t *testing.T) {
	cases := []struct {
		name         string
		dbError      error
		expectedCode int
	}{
		{
			name:         "Name already exists",
			dbError:      postgres.ErrNameAlreadyExists,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "Other error",
			dbError:      errors.New("some unknown error"),
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			gameModel := fixtures.NewGameModelMock(ctrl)
			userModel := fixtures.NewUserModelMock(ctrl)
			authenticator := fixtures.NewAuthenticatorMock(ctrl)
			srv := newServer(t, &Options{
				Authenticator: authenticator,
				GameModel:     gameModel,
				UserModel:     userModel,
			})
			authenticator.EXPECT().
				DecodeToken(gomock.Eq(token)).
				Return(user, nil)
			actualGame := &models.Game{Name: "ACIII"}
			gameModel.
				EXPECT().
				Insert(actualGame).
				Return(nil, c.dbError)
			userModel.
				EXPECT().
				GetUserByToken(token).
				Return(user, nil)

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, "/games", fixtures.Marshal(t, actualGame))
			r.Header.Set(models.XAuthToken, token)
			srv.ServeHTTP(w, r)

			gassert.StatusCode(t, w, c.expectedCode)
		})
	}
}
