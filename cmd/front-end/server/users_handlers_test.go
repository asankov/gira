package server_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/gira-games/client/pkg/client"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"

	gassert "github.com/asankov/gira/internal/fixtures/assert"

	"github.com/asankov/gira/internal/fixtures"
	"github.com/golang/mock/gomock"
)

var (
	email    = "test@mail.com"
	password = "pass123"
	cookie   = &http.Cookie{
		Name:  "token",
		Value: token,
		Path:  "/",
	}
)

func TestSignupForm(t *testing.T) {
	testFormAt(t, "/users/signup")
}

func TestLoginForm(t *testing.T) {
	testFormAt(t, "/users/login")
}

func testFormAt(t *testing.T, path string) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	renderer := fixtures.NewRendererMock(ctrl)
	srv := newServer(nil, renderer)

	renderer.EXPECT().
		Render(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, path, nil)
	srv.ServeHTTP(w, r)

	gassert.StatusOK(t, w)
}

func TestUserLogin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	apiClientMock := fixtures.NewAPIClientMock(ctrl)
	srv := newServer(apiClientMock, nil)

	apiClientMock.EXPECT().
		LoginUser(gomock.AssignableToTypeOf(ctxType), &client.LoginUserRequest{
			Email:    email,
			Password: password,
		}).
		Return(&client.UserLoginResponse{Token: token}, nil)

	w := httptest.NewRecorder()
	form := url.Values{}
	form.Add("email", email)
	form.Add("password", password)
	r := httptest.NewRequest(http.MethodPost, "/users/login", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	srv.ServeHTTP(w, r)

	gassert.Redirect(t, w, "/")
	cookies := w.Result().Cookies()

	require.Equal(t, 1, len(cookies))

	gotCookie := cookies[0]
	assert.Equal(t, cookie.Name, gotCookie.Name)
	assert.Equal(t, cookie.Value, gotCookie.Value)
	assert.Equal(t, cookie.Path, gotCookie.Path)
}

func TestUserLoginFormError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	apiClientMock := fixtures.NewAPIClientMock(ctrl)
	srv := newServer(apiClientMock, nil)

	apiClientMock.EXPECT().
		LoginUser(gomock.AssignableToTypeOf(ctxType), gomock.Eq(&client.LoginUserRequest{
			Email:    email,
			Password: password,
		})).
		Return(nil, errors.New("error while logging in user"))

	w := httptest.NewRecorder()
	form := url.Values{}
	form.Add("email", email)
	form.Add("password", password)
	r := httptest.NewRequest(http.MethodPost, "/users/login", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	srv.ServeHTTP(w, r)

	gassert.StatusCode(t, w, http.StatusBadRequest)
}

func TestUserLoginClientError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	srv := newServer(nil, nil)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/users/login", nil)
	r.Body = nil
	srv.ServeHTTP(w, r)

	gassert.StatusCode(t, w, http.StatusBadRequest)
}

func TestUserLogout(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	apiClientMock := fixtures.NewAPIClientMock(ctrl)
	srv := newServer(apiClientMock, nil)

	apiClientMock.EXPECT().
		LogoutUser(gomock.AssignableToTypeOf(ctxType), &client.LogoutUserRequest{Token: token}).
		Return(nil)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/users/logout", nil)
	r.AddCookie(&http.Cookie{
		Name:  "token",
		Value: token,
	})
	srv.ServeHTTP(w, r)

	gassert.Redirect(t, w, "/")
}

func TestUserSignup(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	apiClientMock := fixtures.NewAPIClientMock(ctrl)
	srv := newServer(apiClientMock, nil)

	apiClientMock.EXPECT().
		CreateUser(gomock.AssignableToTypeOf(ctxType), &client.CreateUserRequest{
			Email:    email,
			Password: password,
		}).
		Return(nil, nil)

	w := httptest.NewRecorder()
	form := url.Values{}
	form.Add("email", email)
	form.Add("password", password)
	r := httptest.NewRequest(http.MethodPost, "/users/create", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	srv.ServeHTTP(w, r)

	gassert.Redirect(t, w, "/")
}

func TestUserSignupNoCredentials(t *testing.T) {
	testCases := []struct {
		Name     string
		Email    string
		Password string
	}{
		{
			Name:     "No email",
			Password: "pass",
		},
		{
			Name:  "No password",
			Email: "test@test.com",
		},
		{
			Name: "No email and password",
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {

			srv := newServer(nil, nil)

			w := httptest.NewRecorder()
			form := url.Values{}
			form.Add("email", testCase.Email)
			form.Add("password", testCase.Password)
			r := httptest.NewRequest(http.MethodPost, "/users/create", strings.NewReader(form.Encode()))
			r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
			srv.ServeHTTP(w, r)

			gassert.StatusCode(t, w, http.StatusBadRequest)
		})
	}
}
