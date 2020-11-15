package server_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
	"testing"

	"github.com/gira-games/client/pkg/client"

	"github.com/asankov/gira/cmd/front-end/server"
	"github.com/asankov/gira/internal/fixtures"
	"github.com/asankov/gira/internal/fixtures/assert"
	"github.com/golangcollege/sessions"
	"github.com/sirupsen/logrus"

	"github.com/golang/mock/gomock"
)

var (
	token = "my-test-token"
	games = []*client.Game{game}
	game  = &client.Game{
		ID:   "1",
		Name: "Game1",
	}
	user = &client.User{
		ID:       "1",
		Username: "test-user",
	}

	ctxType = reflect.TypeOf((*context.Context)(nil)).Elem()
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
		Render(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil)

	srv.ServeHTTP(w, r)

	assert.StatusOK(t, w)
}

func TestHandleHomeLoggedInUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	renderer := fixtures.NewRendererMock(ctrl)
	apiClient := fixtures.NewAPIClientMock(ctrl)

	srv := newServer(apiClient, renderer)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.AddCookie(&http.Cookie{
		Name:  "token",
		Value: token,
	})

	apiClient.EXPECT().
		GetUser(gomock.AssignableToTypeOf(ctxType), &client.GetUserRequest{
			Token: token,
		}).
		Return(&client.GetUserResponse{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
		}, nil)
	renderer.EXPECT().
		Render(gomock.Any(), gomock.Any(), gomock.Eq(server.TemplateData{
			User: user,
		}), gomock.Any()).
		Return(nil)

	srv.ServeHTTP(w, r)

	assert.StatusOK(t, w)
}

func TestHandleHomeRendererError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	renderer := fixtures.NewRendererMock(ctrl)

	srv := newServer(nil, renderer)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)

	renderer.EXPECT().
		Render(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(errors.New("error while rendering page"))

	srv.ServeHTTP(w, r)

	assert.StatusCode(t, w, http.StatusInternalServerError)
}

func TestHandleCreateView(t *testing.T) {
	testCases := []struct {
		Name            string
		Franchises      []*client.Franchise
		FranchisesError error
	}{
		{
			Name:            "GetFranchises returns empty array of franchises and no error",
			Franchises:      []*client.Franchise{},
			FranchisesError: nil,
		},
		{
			Name: "GetFranchises returns array of franchises and no error",
			Franchises: []*client.Franchise{
				{
					ID:   "1",
					Name: "Batman",
				},
			},
			FranchisesError: nil,
		},
		{
			Name:            "GetFranchises returns error",
			Franchises:      nil,
			FranchisesError: errors.New("GetFranchises error"),
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			rendererMock := fixtures.NewRendererMock(ctrl)
			apiClientMock := fixtures.NewAPIClientMock(ctrl)

			srv := newServer(apiClientMock, rendererMock)

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/games/new", nil)
			r.AddCookie(&http.Cookie{
				Name:  "token",
				Value: token,
			})

			apiClientMock.EXPECT().
				GetUser(gomock.AssignableToTypeOf(ctxType), &client.GetUserRequest{Token: token}).
				Return(&client.GetUserResponse{
					ID:       user.ID,
					Username: user.Username,
					Email:    user.Email,
				}, nil)
			apiClientMock.EXPECT().
				GetFranchises(gomock.AssignableToTypeOf(ctxType), &client.GetFranchisesRequest{Token: token}).
				Return(&client.GetFranchisesResponse{Franchises: testCase.Franchises}, testCase.FranchisesError)

			rendererMock.EXPECT().
				Render(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
				Return(nil)

			srv.ServeHTTP(w, r)

			assert.StatusOK(t, w)
		})
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
					GetGames(gomock.AssignableToTypeOf(ctxType), &client.GetGamesRequest{Token: token, ExcludeAssigned: true}).
					Return(&client.GetGamesResponse{Games: games}, nil)
				a.EXPECT().
					GetUser(gomock.AssignableToTypeOf(ctxType), &client.GetUserRequest{Token: token}).
					Return(&client.GetUserResponse{
						ID:       user.ID,
						Username: user.Username,
						Email:    user.Email,
					}, nil)
				r.EXPECT().
					Render(gomock.Any(), gomock.Any(), gomock.Eq(server.TemplateData{Games: games, User: user}), gomock.Any()).
					Return(nil)
			},
		},
		{
			name: "User is not fetched succesfully from the API",
			setup: func(a *fixtures.APIClientMock, r *fixtures.RendererMock) {
				a.EXPECT().
					GetGames(gomock.AssignableToTypeOf(ctxType), &client.GetGamesRequest{Token: token, ExcludeAssigned: true}).
					Return(&client.GetGamesResponse{Games: games}, nil)
				a.EXPECT().
					GetUser(gomock.AssignableToTypeOf(ctxType), &client.GetUserRequest{Token: token}).
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
					GetGames(gomock.AssignableToTypeOf(ctxType), &client.GetGamesRequest{Token: token, ExcludeAssigned: true}).
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
					GetGames(gomock.AssignableToTypeOf(ctxType), gomock.Eq(&client.GetGamesRequest{Token: token, ExcludeAssigned: true})).
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
		LinkGameToUser(gomock.AssignableToTypeOf(ctxType), &client.LinkGameToUserRequest{
			Token:  token,
			GameID: game.ID,
		}).
		Return(nil)

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
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	srv := newServer(nil, nil)

	w := httptest.NewRecorder()
	form := url.Values{}
	r := httptest.NewRequest(http.MethodPost, "/games/add", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.AddCookie(&http.Cookie{
		Name:  "token",
		Value: token,
	})

	srv.ServeHTTP(w, r)

	assert.Redirect(t, w, "/games/add")
}

func TestGamesAddPostClientError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	apiClientMock := fixtures.NewAPIClientMock(ctrl)

	srv := newServer(apiClientMock, nil)

	apiClientMock.EXPECT().
		LinkGameToUser(gomock.AssignableToTypeOf(ctxType), &client.LinkGameToUserRequest{
			Token:  token,
			GameID: game.ID,
		}).
		Return(errors.New("error while linking game"))

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
		UpdateGameProgress(gomock.AssignableToTypeOf(ctxType), &client.UpdateGameProgressRequest{
			GameID: game.ID,
			Token:  token,
			Update: client.UpdateGameProgressChange{
				Status: client.Status("In Progress"),
			},
		}).
		Return(nil)

	w := httptest.NewRecorder()

	form := url.Values{}
	form.Add("game", game.ID)
	form.Add("status", "In Progress")
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
		UpdateGameProgress(gomock.AssignableToTypeOf(ctxType), &client.UpdateGameProgressRequest{
			GameID: game.ID,
			Token:  token,
			Update: client.UpdateGameProgressChange{
				Status: client.Status("In Progress"),
			},
		}).
		Return(errors.New("error while changing status"))

	w := httptest.NewRecorder()

	form := url.Values{}
	form.Add("game", game.ID)
	form.Add("status", "In Progress")
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
					"status": []string{"TODO"},
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

func TestGamesChangeProgress(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	apiClientMock := fixtures.NewAPIClientMock(ctrl)

	srv := newServer(apiClientMock, nil)

	apiClientMock.EXPECT().
		UpdateGameProgress(gomock.AssignableToTypeOf(ctxType), &client.UpdateGameProgressRequest{
			GameID: game.ID,
			Token:  token,
			Update: client.UpdateGameProgressChange{
				Progress: &client.UserGameProgress{
					Current: 10,
					Final:   100,
				},
			},
		}).
		Return(nil)

	w := httptest.NewRecorder()

	form := url.Values{}
	form.Add("game", game.ID)
	form.Add("currentProgress", fmt.Sprintf("%d", 10))
	form.Add("finalProgress", fmt.Sprintf("%d", 100))
	r := httptest.NewRequest(http.MethodPost, "/games/progress", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.AddCookie(&http.Cookie{
		Name:  "token",
		Value: token,
	})
	srv.ServeHTTP(w, r)

	assert.Redirect(t, w, "/games")
}

func TestGamesChangeProgressPostError(t *testing.T) {
	testCases := []struct {
		Name    string
		Current string
		Final   string
		GameID  string
	}{
		{
			Name:    "Game is missing",
			Current: "10",
			Final:   "100",
		},
		{
			Name:    "Current is missing",
			Current: "",
			Final:   "100",
			GameID:  "1",
		},
		{
			Name:    "Final is missing",
			Current: "10",
			Final:   "",
			GameID:  "1",
		},
		{
			Name:    "Current is not a valid int",
			Current: "invalid",
			Final:   "100",
			GameID:  "1",
		},
		{
			Name:    "Final is not a valid int",
			Current: "10",
			Final:   "invalid",
			GameID:  "1",
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			srv := newServer(nil, nil)

			w := httptest.NewRecorder()

			form := url.Values{}
			form.Add("game", testCase.GameID)
			form.Add("currentProgress", testCase.Current)
			form.Add("finalProgress", testCase.Final)
			r := httptest.NewRequest(http.MethodPost, "/games/progress", strings.NewReader(form.Encode()))
			r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
			r.AddCookie(&http.Cookie{
				Name:  "token",
				Value: token,
			})
			srv.ServeHTTP(w, r)

			assert.StatusCode(t, w, http.StatusBadRequest)
		})
	}
}

func TestGamesCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	apiClientMock := fixtures.NewAPIClientMock(ctrl)

	srv := newServer(apiClientMock, nil)

	apiClientMock.EXPECT().
		CreateGame(gomock.AssignableToTypeOf(ctxType), &client.CreateGameRequest{
			Token: token,
			Game: &client.Game{
				Name: game.Name,
			},
		}).
		Return(&client.CreateGameResponse{
			Game: game,
		}, nil)

	w := httptest.NewRecorder()

	form := url.Values{}
	form.Add("name", game.Name)
	r := httptest.NewRequest(http.MethodPost, "/games", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.AddCookie(&http.Cookie{
		Name:  "token",
		Value: token,
	})
	srv.ServeHTTP(w, r)

	assert.Redirect(t, w, "/games/add")
}

func TestGamesCreatePostError(t *testing.T) {
	testCases := []struct {
		name    string
		request func() *http.Request
	}{
		{
			name: "Empty game",
			request: func() *http.Request {
				form := url.Values{}
				body := strings.NewReader(form.Encode())
				r := httptest.NewRequest(http.MethodPost, "/games", body)
				r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

				return r
			},
		},
		{
			name: "Nil form",
			request: func() *http.Request {
				r := httptest.NewRequest(http.MethodPost, "/games", nil)
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

func TestGamesCreateServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	apiClientMock := fixtures.NewAPIClientMock(ctrl)

	srv := newServer(apiClientMock, nil)

	apiClientMock.EXPECT().
		CreateGame(gomock.AssignableToTypeOf(ctxType), &client.CreateGameRequest{Game: &client.Game{Name: game.Name}, Token: token}).
		Return(nil, errors.New("error while creating game"))

	w := httptest.NewRecorder()

	form := url.Values{}
	form.Add("name", game.Name)
	r := httptest.NewRequest(http.MethodPost, "/games", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.AddCookie(&http.Cookie{
		Name:  "token",
		Value: token,
	})
	srv.ServeHTTP(w, r)

	assert.Redirect(t, w, "/games/new")
}

func TestGamesGet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	rendererMock := fixtures.NewRendererMock(ctrl)
	apiClientMock := fixtures.NewAPIClientMock(ctrl)
	srv := newServer(apiClientMock, rendererMock)

	apiClientMock.EXPECT().
		GetUser(gomock.AssignableToTypeOf(ctxType), &client.GetUserRequest{Token: token}).
		Return(&client.GetUserResponse{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
		}, nil)
	apiClientMock.EXPECT().
		GetUserGames(gomock.AssignableToTypeOf(ctxType), &client.GetUserGamesRequest{Token: token}).
		Return(&client.GetUserGamesResponse{
			UserGames: map[client.Status][]*client.UserGame{
				"Done": {
					&client.UserGame{
						ID: "1",
						Game: &client.Game{
							ID:   "1",
							Name: "1",
						},
					},
				},
				"TODO": {
					&client.UserGame{
						ID: "2",
						Game: &client.Game{
							ID:   "2",
							Name: "2",
						},
					},
				},
			},
		}, nil)
	apiClientMock.EXPECT().GetStatuses(gomock.AssignableToTypeOf(ctxType), &client.GetStatusesRequest{Token: token}).Return(&client.GetStatusesResponse{
		Statuses: []client.Status{
			client.Status("TODO"),
			client.Status("Done"),
		},
	}, nil)

	rendererMock.EXPECT().
		Render(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil)

	w := httptest.NewRecorder()

	form := url.Values{}
	form.Add("name", game.Name)
	r := httptest.NewRequest(http.MethodGet, "/games", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.AddCookie(&http.Cookie{
		Name:  "token",
		Value: token,
	})

	srv.ServeHTTP(w, r)
}

func TestGamesGetClientError(t *testing.T) {
	testCases := []struct {
		name   string
		setup  func(*fixtures.APIClientMock)
		assert func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name: "Auth error",
			setup: func(a *fixtures.APIClientMock) {
				a.EXPECT().
					GetUserGames(gomock.AssignableToTypeOf(ctxType), &client.GetUserGamesRequest{Token: token}).
					Return(nil, client.ErrNoAuthorization)
			},

			assert: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert.Redirect(t, w, "/users/login")
			},
		},
		{
			name: "Other error",
			setup: func(a *fixtures.APIClientMock) {
				a.EXPECT().
					GetUserGames(gomock.AssignableToTypeOf(ctxType), &client.GetUserGamesRequest{Token: token}).
					Return(nil, errors.New("unknown error"))
			},
			assert: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert.StatusCode(t, w, http.StatusInternalServerError)
			},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			apiClientMock := fixtures.NewAPIClientMock(ctrl)
			srv := newServer(apiClientMock, nil)

			testCase.setup(apiClientMock)

			w := httptest.NewRecorder()

			form := url.Values{}
			form.Add("name", game.Name)
			r := httptest.NewRequest(http.MethodGet, "/games", strings.NewReader(form.Encode()))
			r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
			r.AddCookie(&http.Cookie{
				Name:  "token",
				Value: token,
			})

			srv.ServeHTTP(w, r)

			testCase.assert(t, w)
		})
	}
}

func TestGamesDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	apiClientMock := fixtures.NewAPIClientMock(ctrl)

	srv := newServer(apiClientMock, nil)

	apiClientMock.EXPECT().
		DeleteUserGame(gomock.AssignableToTypeOf(ctxType), &client.DeleteUserGameRequest{
			GameID: game.ID,
			Token:  token,
		}).
		Return(nil)

	w := httptest.NewRecorder()

	form := url.Values{}
	form.Add("game", game.ID)
	r := httptest.NewRequest(http.MethodPost, "/games/delete", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.AddCookie(&http.Cookie{
		Name:  "token",
		Value: token,
	})
	srv.ServeHTTP(w, r)

	assert.Redirect(t, w, "/games")
}

func TestGamesDeleteClientError(t *testing.T) {
	testCases := []struct {
		name   string
		setup  func(*fixtures.APIClientMock)
		assert func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name: "Auth error",
			setup: func(a *fixtures.APIClientMock) {
				a.EXPECT().
					DeleteUserGame(gomock.AssignableToTypeOf(ctxType), &client.DeleteUserGameRequest{
						GameID: game.ID,
						Token:  token,
					}).
					Return(client.ErrNoAuthorization)
			},

			assert: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert.Redirect(t, w, "/users/login")
			},
		},
		{
			name: "Other error",
			setup: func(a *fixtures.APIClientMock) {
				a.EXPECT().
					DeleteUserGame(gomock.AssignableToTypeOf(ctxType), &client.DeleteUserGameRequest{
						GameID: game.ID,
						Token:  token,
					}).
					Return(errors.New("unknown error"))

			},
			assert: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert.StatusCode(t, w, http.StatusInternalServerError)
			},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			apiClientMock := fixtures.NewAPIClientMock(ctrl)
			srv := newServer(apiClientMock, nil)

			testCase.setup(apiClientMock)

			w := httptest.NewRecorder()

			form := url.Values{}
			form.Add("game", game.ID)
			r := httptest.NewRequest(http.MethodPost, "/games/delete", strings.NewReader(form.Encode()))
			r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
			r.AddCookie(&http.Cookie{
				Name:  "token",
				Value: token,
			})

			srv.ServeHTTP(w, r)

			testCase.assert(t, w)
		})
	}
}

func TestGamesDeletePostError(t *testing.T) {
	testCases := []struct {
		name    string
		request func() *http.Request
	}{
		{
			name: "Empty game",
			request: func() *http.Request {
				form := url.Values{
					"status": []string{"TODO"},
				}
				body := strings.NewReader(form.Encode())
				r := httptest.NewRequest(http.MethodPost, "/games/delete", body)
				r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

				return r
			},
		},
		{
			name: "Nil form",
			request: func() *http.Request {
				r := httptest.NewRequest(http.MethodPost, "/games/delete", nil)
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
