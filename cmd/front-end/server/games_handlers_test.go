package server_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/asankov/gira/cmd/front-end/server"
	"github.com/asankov/gira/internal/fixtures"
	"github.com/asankov/gira/pkg/client"
	"github.com/asankov/gira/pkg/models"
	"github.com/golangcollege/sessions"
	"github.com/sirupsen/logrus"

	"github.com/golang/mock/gomock"
)

var (
	token = "my-test-token"
	games = []*models.Game{game}
	game  = &models.Game{
		ID:   "1",
		Name: "Game1",
	}
	user = &models.User{
		ID:       "1",
		Username: "test-user",
	}
)

func newServer(a *fixtures.APIClientMock, r *fixtures.RendererMock) *server.Server {
	session := sessions.New([]byte("secret"))

	return &server.Server{
		Log:      logrus.StandardLogger(),
		Renderer: r,
		Client:   a,
		Session:  session,
	}
}

func TestHandleHome(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	renderer := fixtures.NewRendererMock(ctrl)

	srv := newServer(nil, renderer)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)

	renderer.EXPECT().
		Render(gomock.Eq(w), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil)

	srv.ServeHTTP(w, r)

	got, expected := w.Code, http.StatusOK
	if got != expected {
		t.Errorf("Got (%d) for status code, expected (%d)", got, expected)
	}
}

func TestGamesAdd(t *testing.T) {
	testCases := []struct {
		name  string
		setup func(*fixtures.APIClientMock, *fixtures.RendererMock)
	}{
		{
			name: "User is fetched succesfully from the API",
			setup: func(a *fixtures.APIClientMock, r *fixtures.RendererMock) {
				a.EXPECT().
					GetGames(gomock.Eq(token), gomock.Any()).
					Return(games, nil)
				a.EXPECT().
					GetUser(gomock.Eq(token)).
					Return(user, nil)
				r.EXPECT().
					Render(gomock.Any(), gomock.Any(), gomock.Eq(server.TemplateData{Games: games, User: user}), gomock.Any()).
					Return(nil)
			},
		},
		{
			name: "User is not fetched succesfully from the API",
			setup: func(a *fixtures.APIClientMock, r *fixtures.RendererMock) {
				a.EXPECT().
					GetGames(gomock.Eq(token), gomock.Any()).
					Return(games, nil)
				a.EXPECT().
					GetUser(gomock.Eq(token)).
					Return(nil, errors.New("error while fetching user"))
				r.EXPECT().
					Render(gomock.Any(), gomock.Any(), gomock.Eq(server.TemplateData{Games: games}), gomock.Any()).
					Return(nil)
			},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			apiClientMock := fixtures.NewAPIClientMock(ctrl)
			rendererMock := fixtures.NewRendererMock(ctrl)

			testCase.setup(apiClientMock, rendererMock)

			srv := newServer(apiClientMock, rendererMock)

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/games/add", nil)
			r.AddCookie(&http.Cookie{
				Name:  "token",
				Value: token,
			})

			srv.ServeHTTP(w, r)

			got, expected := w.Code, http.StatusOK
			if got != expected {
				t.Errorf("Got (%d) for status code, expected (%d)", got, expected)
			}
		})
	}
}

func TestGamesAddClientError(t *testing.T) {
	testCases := []struct {
		name              string
		setup             func(*fixtures.APIClientMock)
		expectedCode      int
		additionalAsserts func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name: "NoAuthorization Error",
			setup: func(a *fixtures.APIClientMock) {
				a.EXPECT().
					GetGames(gomock.Eq(token), gomock.Any()).
					Return(nil, client.ErrNoAuthorization)
			},

			expectedCode: http.StatusSeeOther,
			additionalAsserts: func(t *testing.T, w *httptest.ResponseRecorder) {
				got, expected := w.Header().Get("Location"), "/users/login"
				if got != expected {
					t.Errorf("Got %s for Location header, expected %s", got, expected)
				}
			},
		},
		{
			name: "Other error",
			setup: func(a *fixtures.APIClientMock) {
				a.EXPECT().
					GetGames(gomock.Eq(token), gomock.Any()).
					Return(nil, errors.New("some other error"))
			},
			expectedCode:      http.StatusInternalServerError,
			additionalAsserts: func(*testing.T, *httptest.ResponseRecorder) {},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			apiClientMock := fixtures.NewAPIClientMock(ctrl)

			testCase.setup(apiClientMock)

			srv := newServer(apiClientMock, nil)

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/games/add", nil)
			r.AddCookie(&http.Cookie{
				Name:  "token",
				Value: token,
			})

			srv.ServeHTTP(w, r)

			got, expected := w.Code, testCase.expectedCode
			if got != expected {
				t.Errorf("Got (%d) for status code, expected (%d)", got, testCase.expectedCode)
			}
			testCase.additionalAsserts(t, w)
		})
	}

}
