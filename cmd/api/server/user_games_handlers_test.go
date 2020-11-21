package server

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	gassert "github.com/asankov/gira/internal/fixtures/assert"

	"github.com/asankov/gira/internal/fixtures"
	"github.com/asankov/gira/pkg/models"
	"github.com/golang/mock/gomock"
)

func TestUsersGamesPatch(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	authenticatorMock := fixtures.NewAuthenticatorMock(ctrl)
	gamesModelMock := fixtures.NewGameModelMock(ctrl)
	userModelMock := fixtures.NewUserModelMock(ctrl)
	srv := newServer(t, &Options{
		Authenticator: authenticatorMock,
		UserModel:     userModelMock,
		GameModel:     gamesModelMock,
	})

	authenticatorMock.EXPECT().
		DecodeToken(gomock.Eq(token)).
		Return(nil, nil)
	userModelMock.EXPECT().
		GetUserByToken(gomock.Eq(token)).
		Return(&models.User{
			ID: "12",
		}, nil)
	gamesModelMock.
		EXPECT().
		ChangeGameStatus(gomock.Eq("12"), gomock.Eq("1"), gomock.Eq(models.StatusDone)).
		Return(nil)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPatch, "/games/1", fixtures.Marshal(t, models.ChangeGameStatusRequest{Status: models.StatusDone}))
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
	r := httptest.NewRequest(http.MethodPatch, "/games/1", nil)
	r.Header.Add(models.XAuthToken, token)

	srv.ServeHTTP(w, r)

	gassert.StatusCode(t, w, http.StatusBadRequest)
}

func TestUsersGamesPatchServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	authenticatorMock := fixtures.NewAuthenticatorMock(ctrl)
	gamesModelMock := fixtures.NewGameModelMock(ctrl)
	userModelMock := fixtures.NewUserModelMock(ctrl)
	srv := newServer(t, &Options{
		Authenticator: authenticatorMock,
		UserModel:     userModelMock,
		GameModel:     gamesModelMock,
	})

	authenticatorMock.EXPECT().
		DecodeToken(gomock.Eq(token)).
		Return(nil, nil)
	userModelMock.EXPECT().
		GetUserByToken(gomock.Eq(token)).
		Return(&models.User{
			ID: "12",
		}, nil)
	gamesModelMock.
		EXPECT().
		ChangeGameStatus(gomock.Eq("12"), gomock.Eq("1"), gomock.Eq(models.StatusDone)).
		Return(errors.New("error while changing game status"))

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPatch, "/games/1", fixtures.Marshal(t, models.ChangeGameStatusRequest{Status: models.StatusDone}))
	r.Header.Add(models.XAuthToken, token)

	srv.ServeHTTP(w, r)

	gassert.StatusCode(t, w, http.StatusInternalServerError)
}

func TestUsersGamesPatchInvalidInput(t *testing.T) {
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
	r := httptest.NewRequest(http.MethodPatch, "/games/1", fixtures.Marshal(t, models.ChangeGameStatusRequest{Status: models.Status("some status")}))
	r.Header.Add(models.XAuthToken, token)

	srv.ServeHTTP(w, r)

	gassert.StatusCode(t, w, http.StatusBadRequest)
}

func TestUsersGamesDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	authenticatorMock := fixtures.NewAuthenticatorMock(ctrl)
	gamesModelMock := fixtures.NewGameModelMock(ctrl)
	userModelMock := fixtures.NewUserModelMock(ctrl)
	srv := newServer(t, &Options{
		Authenticator: authenticatorMock,
		UserModel:     userModelMock,
		GameModel:     gamesModelMock,
	})

	authenticatorMock.EXPECT().
		DecodeToken(gomock.Eq(token)).
		Return(nil, nil)
	userModelMock.EXPECT().
		GetUserByToken(gomock.Eq(token)).
		Return(&models.User{
			ID: "12",
		}, nil)
	gamesModelMock.
		EXPECT().
		DeleteGame(gomock.Eq("12"), gomock.Eq("1")).
		Return(nil)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodDelete, "/games/1", nil)
	r.Header.Add(models.XAuthToken, token)

	srv.ServeHTTP(w, r)

	gassert.StatusOK(t, w)
}

func TestUsersGamesDeleteUserDoesNotOwnGame(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	authenticatorMock := fixtures.NewAuthenticatorMock(ctrl)
	gamesModelMock := fixtures.NewGameModelMock(ctrl)
	userModelMock := fixtures.NewUserModelMock(ctrl)
	srv := newServer(t, &Options{
		Authenticator: authenticatorMock,
		UserModel:     userModelMock,
		GameModel:     gamesModelMock,
	})

	authenticatorMock.EXPECT().
		DecodeToken(gomock.Eq(token)).
		Return(nil, nil)
	userModelMock.EXPECT().
		GetUserByToken(gomock.Eq(token)).
		Return(&models.User{
			ID: "12",
		}, nil)
	gamesModelMock.EXPECT().
		DeleteGame(gomock.Eq("12"), gomock.Eq("1")).
		Return(fmt.Errorf("user does not own game"))

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodDelete, "/games/1", nil)
	r.Header.Add(models.XAuthToken, token)

	srv.ServeHTTP(w, r)

	gassert.StatusCode(t, w, http.StatusBadRequest)
}
