package client_test

import (
	"net/http"
	"testing"

	"github.com/gira-games/api/internal/fixtures"
	"github.com/gira-games/api/pkg/client"
	"github.com/gira-games/api/pkg/models"
	"github.com/stretchr/testify/require"
)

var (
	franchise = &models.Franchise{
		ID:   "1",
		Name: "Batman",
	}
	franchises         = []*models.Franchise{franchise}
	franchisesResponse = &models.FranchisesResponse{
		Franchises: franchises,
	}
)

func TestFranchisesGet(t *testing.T) {
	ts := fixtures.NewTestServer(t).
		Path("/franchises").
		Token(token).
		Method(http.MethodGet).
		Data(franchisesResponse).
		Build()
	defer ts.Close()

	cl := newClient(t, ts.URL)

	fr, err := cl.GetFranchises(token)
	require.NoError(t, err)
	require.Equal(t, franchises, fr)
}

func TestFranchisesCreate(t *testing.T) {
	ts := fixtures.NewTestServer(t).
		Path("/franchises").
		Token(token).
		Method(http.MethodPost).
		Data(franchise).
		Build()
	defer ts.Close()

	cl := newClient(t, ts.URL)

	frResponse, err := cl.CreateFranchise(&client.CreateFranchiseRequest{
		Name: "AC",
	}, token)
	require.NoError(t, err)
	require.Equal(t, franchise, frResponse.Franchise)
}
