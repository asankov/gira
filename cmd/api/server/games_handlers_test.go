package server

import (
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/asankov/gira/internal/fixtures"
	"github.com/asankov/gira/pkg/models"
	"github.com/asankov/gira/pkg/models/postgres"
	"github.com/golang/mock/gomock"
)

var (
	user = &models.User{Username: "anton"}
)

func setupGamesServer(g GameModel, a *fixtures.AuthenticatorMock) *Server {
	a.EXPECT().
		DecodeToken(gomock.Eq(token)).
		Return(user, nil)

	return &Server{
		Log:           log.New(os.Stdout, "", 0),
		GameModel:     g,
		Authenticator: a,
	}
}

func TestGetGames(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	gameModel := fixtures.NewGameModelMock(ctrl)
	authenticator := fixtures.NewAuthenticatorMock(ctrl)
	srv := setupGamesServer(gameModel, authenticator)

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
	r.Header.Set("x-auth-token", token)
	srv.ServeHTTP(w, r)

	got, expected := w.Result().StatusCode, http.StatusOK
	if got != expected {
		t.Fatalf("Got status code - (%d), expected (%d)", got, expected)
	}

	var res []*models.Game
	fixtures.Decode(t, w.Body, &res)

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
	defer ctrl.Finish()

	gameModel := fixtures.NewGameModelMock(ctrl)
	authenticator := fixtures.NewAuthenticatorMock(ctrl)
	srv := setupGamesServer(gameModel, authenticator)

	gameModel.
		EXPECT().
		All().
		Return(nil, errors.New("this is an intentional error"))

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/games", nil)
	r.Header.Set("x-auth-token", token)
	srv.ServeHTTP(w, r)

	got, expected := w.Result().StatusCode, http.StatusInternalServerError
	if got != expected {
		t.Fatalf("Got status code - (%d), expected (%d)", got, expected)
	}
}

func TestGetGameByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	gameModel := fixtures.NewGameModelMock(ctrl)
	authenticator := fixtures.NewAuthenticatorMock(ctrl)
	srv := setupGamesServer(gameModel, authenticator)

	actualName := "ACIII"
	actualGame := &models.Game{Name: actualName}
	gameModel.
		EXPECT().
		Get("1").
		Return(actualGame, nil)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/games/1", nil)
	r.Header.Set("x-auth-token", token)
	srv.ServeHTTP(w, r)

	statusCode := w.Result().StatusCode
	if statusCode != http.StatusOK {
		t.Fatalf("Got (%d) for status code, expected (%d)", statusCode, http.StatusOK)
	}

	var game *models.Game
	fixtures.Decode(t, w.Body, &game)
	if game.Name != actualName {
		t.Fatalf("Got (%s) for game.Name, expected (%s)", game.Name, actualName)
	}
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
			authenticator := fixtures.NewAuthenticatorMock(ctrl)
			srv := setupGamesServer(gameModel, authenticator)

			gameModel.
				EXPECT().
				Get("1").
				Return(nil, c.dbError)

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/games/1", nil)
			r.Header.Set("x-auth-token", token)
			srv.ServeHTTP(w, r)

			got, expected := w.Result().StatusCode, c.expectedCode
			if got != expected {
				t.Fatalf("Got (%d) for status code, expected (%d)", got, expected)
			}
		})
	}

}
func TestCreateGame(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	gameModel := fixtures.NewGameModelMock(ctrl)
	authenticator := fixtures.NewAuthenticatorMock(ctrl)

	srv := setupGamesServer(gameModel, authenticator)

	actualName := "ACIII"
	actualGame := &models.Game{Name: actualName}
	gameModel.
		EXPECT().
		Insert(actualGame).
		Return(actualGame, nil)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/games", fixtures.Marshall(t, actualGame))
	r.Header.Set("x-auth-token", token)
	srv.ServeHTTP(w, r)

	statusCode := w.Result().StatusCode
	if statusCode != http.StatusOK {
		t.Fatalf("Got (%d) for status code, expected (%d)", statusCode, http.StatusOK)
	}

	var game *models.Game
	fixtures.Decode(t, w.Body, &game)
	if game.Name != actualName {
		t.Fatalf("Got (%s) for game.Name, expected (%s)", game.Name, actualName)
	}
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
			authenticator := fixtures.NewAuthenticatorMock(ctrl)

			srv := setupGamesServer(gameModel, authenticator)

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, "/games", fixtures.Marshall(t, c.game))
			r.Header.Set("x-auth-token", token)
			srv.ServeHTTP(w, r)

			got, expected := w.Result().StatusCode, http.StatusBadRequest
			if got != expected {
				t.Fatalf("Got (%d) for status code, expected (%d)", got, expected)
			}
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
			authenticator := fixtures.NewAuthenticatorMock(ctrl)
			srv := setupGamesServer(gameModel, authenticator)

			actualGame := &models.Game{Name: "ACIII"}
			gameModel.
				EXPECT().
				Insert(actualGame).
				Return(nil, c.dbError)

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, "/games", fixtures.Marshall(t, actualGame))
			r.Header.Set("x-auth-token", token)
			srv.ServeHTTP(w, r)

			got, expected := w.Result().StatusCode, c.expectedCode
			if got != expected {
				t.Fatalf("Got (%d) for status code, expected (%d)", got, expected)
			}
		})
	}
}
