package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/asankov/gira/pkg/models/postgres"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	gassert "github.com/asankov/gira/internal/fixtures/assert"

	"github.com/asankov/gira/pkg/models"

	"github.com/asankov/gira/internal/fixtures"
	"github.com/golang/mock/gomock"
)

var (
	franchiseBatman = models.Franchise{
		ID:   "123",
		Name: "Batman",
	}
	franchiseAC = models.Franchise{
		ID:   "124",
		Name: "AC",
	}
	franchises = []*models.Franchise{&franchiseBatman, &franchiseAC}
)

func TestFranchisesGet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	franchiseModel := fixtures.NewFranchiseModelMock(ctrl)
	userModel := fixtures.NewUserModelMock(ctrl)
	authenticator := fixtures.NewAuthenticatorMock(ctrl)
	srv := newServer(t, &Options{
		Authenticator:  authenticator,
		UserModel:      userModel,
		FranchiseModel: franchiseModel,
	})

	authenticator.EXPECT().
		DecodeToken(gomock.Eq(token)).
		Return(user, nil)
	userModel.
		EXPECT().
		GetUserByToken(token).
		Return(user, nil)
	franchiseModel.EXPECT().
		All().
		Return(franchises, nil)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/franchises", nil)
	r.Header.Set(models.XAuthToken, token)
	srv.ServeHTTP(w, r)

	var res models.FranchisesResponse
	fixtures.Decode(t, w.Body, &res)

	gassert.StatusOK(t, w)
	require.Equal(t, len(franchises), len(res.Franchises))
	for _, fr := range franchises {
		assert.Contains(t, res.Franchises, fr)
	}
}

func TestFranchisesCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	franchiseModel := fixtures.NewFranchiseModelMock(ctrl)
	userModel := fixtures.NewUserModelMock(ctrl)
	authenticator := fixtures.NewAuthenticatorMock(ctrl)
	srv := newServer(t, &Options{
		Authenticator:  authenticator,
		UserModel:      userModel,
		FranchiseModel: franchiseModel,
	})

	authenticator.EXPECT().
		DecodeToken(gomock.Eq(token)).
		Return(user, nil)
	userModel.
		EXPECT().
		GetUserByToken(token).
		Return(user, nil)
	franchiseModel.EXPECT().
		Insert(&franchiseBatman).
		Return(&franchiseBatman, nil)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/franchises", fixtures.Marshal(t, franchiseBatman))
	r.Header.Set(models.XAuthToken, token)
	srv.ServeHTTP(w, r)

	var res models.Franchise
	fixtures.Decode(t, w.Body, &res)

	gassert.StatusOK(t, w)
	require.Equal(t, franchiseBatman, res)
}

func TestFranchisesCreateValidationError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userModel := fixtures.NewUserModelMock(ctrl)
	authenticator := fixtures.NewAuthenticatorMock(ctrl)
	srv := newServer(t, &Options{
		Authenticator: authenticator,
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
	r := httptest.NewRequest(http.MethodPost, "/franchises", fixtures.Marshal(t, &models.Franchise{}))
	r.Header.Set(models.XAuthToken, token)
	srv.ServeHTTP(w, r)

	gassert.StatusCode(t, w, http.StatusBadRequest)
	// TODO: assert body once we start returning JSON
	// require.Equal(t, franchiseBatman, res)
}

func TestFranchisesCreateDBError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	franchiseModel := fixtures.NewFranchiseModelMock(ctrl)
	userModel := fixtures.NewUserModelMock(ctrl)
	authenticator := fixtures.NewAuthenticatorMock(ctrl)
	srv := newServer(t, &Options{
		Authenticator:  authenticator,
		UserModel:      userModel,
		FranchiseModel: franchiseModel,
	})

	authenticator.EXPECT().
		DecodeToken(gomock.Eq(token)).
		Return(user, nil)
	userModel.
		EXPECT().
		GetUserByToken(token).
		Return(user, nil)
	franchiseModel.EXPECT().
		Insert(&franchiseBatman).
		Return(nil, postgres.ErrNameAlreadyExists)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/franchises", fixtures.Marshal(t, franchiseBatman))
	r.Header.Set(models.XAuthToken, token)
	srv.ServeHTTP(w, r)

	gassert.StatusCode(t, w, http.StatusBadRequest)
	// TODO: assert body once we start returning JSON
	// require.Equal(t, franchiseBatman, res)
}
