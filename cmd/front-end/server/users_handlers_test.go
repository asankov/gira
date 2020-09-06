package server_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/asankov/gira/internal/fixtures"
	"github.com/asankov/gira/pkg/models"
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

	got, expected := w.Code, http.StatusOK
	if got != expected {
		t.Errorf("Got (%d) for status code, expected (%d)", got, expected)
	}
}

func TestUserLogin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	apiClientMock := fixtures.NewAPIClientMock(ctrl)
	srv := newServer(apiClientMock, nil)

	apiClientMock.EXPECT().
		LoginUser(gomock.Eq(&models.User{
			Email:    email,
			Password: password,
		})).
		Return(&models.UserLoginResponse{Token: token}, nil)

	w := httptest.NewRecorder()
	form := url.Values{}
	form.Add("email", email)
	form.Add("password", password)
	r := httptest.NewRequest(http.MethodPost, "/users/login", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	srv.ServeHTTP(w, r)

	gotCode, expectedCode := w.Code, http.StatusSeeOther
	if gotCode != expectedCode {
		t.Errorf("Got (%d) for status code, expected (%d)", gotCode, expectedCode)
	}
	gotHeader, expectedHeader := w.Header().Get("Location"), "/"
	if gotHeader != expectedHeader {
		t.Errorf("Got %s for Location header, expected %s", gotHeader, expectedHeader)
	}
	cookies := w.Result().Cookies()
	if len(cookies) != 1 {
		t.Fatalf("Got (%d) cookies, expected 1", len(cookies))
	}
	gotCookie := cookies[0]
	if gotCookie.Name != cookie.Name || gotCookie.Value != cookie.Value || gotCookie.Path != cookie.Path {
		t.Errorf("Got (%v) cookie, expected (%v)", gotCookie, cookie)
	}
}

func TestUserLoginFormError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	apiClientMock := fixtures.NewAPIClientMock(ctrl)
	srv := newServer(apiClientMock, nil)

	apiClientMock.EXPECT().
		LoginUser(gomock.Eq(&models.User{
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

	gotCode, expectedCode := w.Code, http.StatusBadRequest
	if gotCode != expectedCode {
		t.Errorf("Got (%d) for status code, expected (%d)", gotCode, expectedCode)
	}
}

func TestUserLoginClientError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	srv := newServer(nil, nil)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/users/login", nil)
	r.Body = nil
	srv.ServeHTTP(w, r)

	gotCode, expectedCode := w.Code, http.StatusBadRequest
	if gotCode != expectedCode {
		t.Errorf("Got (%d) for status code, expected (%d)", gotCode, expectedCode)
	}
}
