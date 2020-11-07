package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gira-games/api/pkg/models/postgres"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	gassert "github.com/gira-games/api/internal/fixtures/assert"

	"github.com/gira-games/api/pkg/models"

	"github.com/gira-games/api/internal/fixtures"
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
		Insert(&models.Franchise{Name: "Batman"}).
		Return(&franchiseBatman, nil)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/franchises", fixtures.Marshal(t, models.Franchise{Name: "Batman"}))
	r.Header.Set(models.XAuthToken, token)
	srv.ServeHTTP(w, r)

	var res models.Franchise
	fixtures.Decode(t, w.Body, &res)

	gassert.StatusOK(t, w)
	require.Equal(t, franchiseBatman, res)
}

func TestFranchisesCreateValidationError(t *testing.T) {
	testCases := []struct {
		name      string
		franchise *models.Franchise
	}{
		{
			name:      "Empty name",
			franchise: &models.Franchise{},
		},
		{
			name: "Filled ID",
			franchise: &models.Franchise{
				ID:   "123",
				Name: "Batman",
			},
		},
		{
			name: "Empty name and filled ID",
			franchise: &models.Franchise{
				ID: "123",
			},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
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
			r := httptest.NewRequest(http.MethodPost, "/franchises", fixtures.Marshal(t, testCase.franchise))
			r.Header.Set(models.XAuthToken, token)
			srv.ServeHTTP(w, r)

			gassert.StatusCode(t, w, http.StatusBadRequest)
			var err models.ErrorResponse
			fixtures.Decode(t, w.Body, &err)
			require.NotEmpty(t, err.Error, "Error returned from server should not be empty")
		})
	}
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
		Insert(&models.Franchise{Name: "Batman"}).
		Return(nil, postgres.ErrNameAlreadyExists)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/franchises", fixtures.Marshal(t, models.Franchise{Name: "Batman"}))
	r.Header.Set(models.XAuthToken, token)
	srv.ServeHTTP(w, r)

	gassert.StatusCode(t, w, http.StatusBadRequest)
	var err models.ErrorResponse
	fixtures.Decode(t, w.Body, &err)
	require.NotEmpty(t, err.Error, "Error returned from server should not be empty")
}
