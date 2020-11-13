package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gira-games/api/internal/fixtures"
	gassert "github.com/gira-games/api/internal/fixtures/assert"
	"github.com/gira-games/api/pkg/models"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestGetStatuses(t *testing.T) {
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
		Return(user, nil)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/statuses", nil)
	r.Header.Add(models.XAuthToken, token)

	srv.ServeHTTP(w, r)

	gassert.StatusOK(t, w)

	var statusesResponse models.StatusesResponse
	fixtures.Decode(t, w.Body, &statusesResponse)

	require.ElementsMatch(t, statusesResponse.Statuses, models.AllStatuses)
}
