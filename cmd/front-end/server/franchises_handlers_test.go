package server_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/gira-games/client/pkg/client"

	"github.com/asankov/gira/internal/fixtures"
	"github.com/asankov/gira/internal/fixtures/assert"
	"github.com/golang/mock/gomock"
)

var franchise = client.Franchise{
	ID:   "123",
	Name: "Batman",
}

func TestAddFranchise(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	apiClientMock := fixtures.NewAPIClientMock(ctrl)

	srv := newServer(apiClientMock, nil)

	apiClientMock.EXPECT().
		CreateFranchise(gomock.AssignableToTypeOf(ctxType), gomock.Eq(&client.CreateFranchiseRequest{Name: franchise.Name, Token: token})).
		Return(&client.CreateFranchiseResponse{
			Franchise: &franchise,
		}, nil)

	w := httptest.NewRecorder()

	form := url.Values{}
	form.Add("franchise", franchise.Name)
	r := httptest.NewRequest(http.MethodPost, "/franchises/add", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.AddCookie(&http.Cookie{
		Name:  "token",
		Value: token,
	})
	srv.ServeHTTP(w, r)

	assert.Redirect(t, w, fmt.Sprintf("/games/new?selectedFranchise=%s", franchise.ID))
}

func TestAddFranchiseEmptyFranchise(t *testing.T) {

	srv := newServer(nil, nil)

	w := httptest.NewRecorder()

	form := url.Values{}
	r := httptest.NewRequest(http.MethodPost, "/franchises/add", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.AddCookie(&http.Cookie{
		Name:  "token",
		Value: token,
	})
	srv.ServeHTTP(w, r)

	assert.StatusCode(t, w, http.StatusBadRequest)
}

func TestAddFranchiseNoAuth(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	apiClientMock := fixtures.NewAPIClientMock(ctrl)

	srv := newServer(apiClientMock, nil)

	apiClientMock.EXPECT().
		CreateFranchise(gomock.AssignableToTypeOf(ctxType), gomock.Eq(&client.CreateFranchiseRequest{Name: franchise.Name, Token: token})).
		Return(nil, client.ErrNoAuthorization)

	w := httptest.NewRecorder()

	form := url.Values{}
	form.Add("franchise", franchise.Name)
	r := httptest.NewRequest(http.MethodPost, "/franchises/add", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.AddCookie(&http.Cookie{
		Name:  "token",
		Value: token,
	})
	srv.ServeHTTP(w, r)

	assert.Redirect(t, w, "/users/login")
}
