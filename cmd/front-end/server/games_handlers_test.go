package server_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/asankov/gira/cmd/front-end/server"
	"github.com/asankov/gira/internal/fixtures"
	"github.com/asankov/gira/internal/fixtures/assert"
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

	assert.StatusOK(t, w)
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

			assert.StatusOK(t, w)
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
				assert.Redirect(t, w, "/users/login")
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

			assert.StatusCode(t, w, testCase.expectedCode)
			testCase.additionalAsserts(t, w)
		})
	}
}

func TestGamesAddPost(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	apiClientMock := fixtures.NewAPIClientMock(ctrl)

	srv := newServer(apiClientMock, nil)

	apiClientMock.EXPECT().
		LinkGameToUser(gomock.Eq(game.ID), token).
		Return(nil, nil)

	w := httptest.NewRecorder()

	form := url.Values{}
	form.Add("game", game.ID)
	r := httptest.NewRequest(http.MethodPost, "/games/add", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.AddCookie(&http.Cookie{
		Name:  "token",
		Value: token,
	})
	srv.ServeHTTP(w, r)

	assert.Redirect(t, w, "/games")
}

func TestGamesAddPostFormError(t *testing.T) {
	testCases := []struct {
		name       string
		getRequest func() *http.Request
	}{
		{
			name: "Error parsing form",
			getRequest: func() *http.Request {
				r := httptest.NewRequest(http.MethodPost, "/games/add", nil)
				r.Body = nil
				r.AddCookie(&http.Cookie{
					Name:  "token",
					Value: token,
				})
				return r
			},
		},
		{
			name: "Validation error",
			getRequest: func() *http.Request {
				form := url.Values{}
				r := httptest.NewRequest(http.MethodPost, "/games/add", strings.NewReader(form.Encode()))
				r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
				r.AddCookie(&http.Cookie{
					Name:  "token",
					Value: token,
				})
				return r
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			srv := newServer(nil, nil)

			w := httptest.NewRecorder()

			srv.ServeHTTP(w, testCase.getRequest())

			assert.StatusCode(t, w, http.StatusBadRequest)
		})
	}
}

func TestGamesAddPostClientError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	apiClientMock := fixtures.NewAPIClientMock(ctrl)

	srv := newServer(apiClientMock, nil)

	apiClientMock.EXPECT().
		LinkGameToUser(gomock.Eq(game.ID), token).
		Return(nil, errors.New("error while linking game"))

	w := httptest.NewRecorder()

	form := url.Values{}
	form.Add("game", game.ID)
	r := httptest.NewRequest(http.MethodPost, "/games/add", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.AddCookie(&http.Cookie{
		Name:  "token",
		Value: token,
	})
	srv.ServeHTTP(w, r)

	assert.StatusCode(t, w, http.StatusInternalServerError)
}

func TestGamesChangeStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	apiClientMock := fixtures.NewAPIClientMock(ctrl)

	srv := newServer(apiClientMock, nil)

	apiClientMock.EXPECT().
		ChangeGameStatus(gomock.Eq(game.ID), gomock.Eq(token), gomock.Eq(models.StatusInProgress)).
		Return(nil)

	w := httptest.NewRecorder()

	form := url.Values{}
	form.Add("game", game.ID)
	form.Add("status", string(models.StatusInProgress))
	r := httptest.NewRequest(http.MethodPost, "/games/status", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.AddCookie(&http.Cookie{
		Name:  "token",
		Value: token,
	})
	srv.ServeHTTP(w, r)

	assert.Redirect(t, w, "/games")
}

func TestGamesChangeStatusServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	apiClientMock := fixtures.NewAPIClientMock(ctrl)

	srv := newServer(apiClientMock, nil)

	apiClientMock.EXPECT().
		ChangeGameStatus(gomock.Eq(game.ID), gomock.Eq(token), gomock.Eq(models.StatusInProgress)).
		Return(errors.New("error while changing status"))

	w := httptest.NewRecorder()

	form := url.Values{}
	form.Add("game", game.ID)
	form.Add("status", string(models.StatusInProgress))
	r := httptest.NewRequest(http.MethodPost, "/games/status", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.AddCookie(&http.Cookie{
		Name:  "token",
		Value: token,
	})
	srv.ServeHTTP(w, r)

	assert.StatusCode(t, w, http.StatusInternalServerError)
}

func TestGamesChangeStatusPostError(t *testing.T) {
	testCases := []struct {
		name    string
		request func() *http.Request
	}{
		{
			name: "Empty game",
			request: func() *http.Request {
				form := url.Values{
					"status": []string{string(models.StatusTODO)},
				}
				body := strings.NewReader(form.Encode())
				r := httptest.NewRequest(http.MethodPost, "/games/status", body)
				r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

				return r
			},
		},
		{
			name: "Empty status",
			request: func() *http.Request {
				form := url.Values{
					"game": []string{"1"},
				}
				body := strings.NewReader(form.Encode())
				r := httptest.NewRequest(http.MethodPost, "/games/status", body)
				r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

				return r
			},
		},
		{
			name: "Nil form",
			request: func() *http.Request {
				r := httptest.NewRequest(http.MethodPost, "/games/status", nil)
				r.Body = nil
				return r
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			srv := newServer(nil, nil)

			w := httptest.NewRecorder()
			r := testCase.request()
			r.AddCookie(&http.Cookie{
				Name:  "token",
				Value: token,
			})
			srv.ServeHTTP(w, r)

			assert.StatusCode(t, w, http.StatusBadRequest)
		})
	}
}
