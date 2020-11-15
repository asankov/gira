package client_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/asankov/gira/internal/fixtures"
	"github.com/asankov/gira/pkg/client"
	"github.com/stretchr/testify/assert"
)

var statuses []client.Status = []client.Status{
	client.Status("TODO"),
	client.Status("In Progress"),
	client.Status("Done"),
}

func TestGetStatuses(t *testing.T) {
	ts := fixtures.NewTestServer(t).
		Path("/statuses").
		Method(http.MethodGet).
		Token(token).
		Data(client.GetStatusesResponse{Statuses: statuses}).
		Build()
	defer ts.Close()

	cl := newClient(t, ts.URL)

	resp, err := cl.GetStatuses(context.Background(), &client.GetStatusesRequest{
		Token: token,
	})
	assert.NoError(t, err)
	assert.EqualValues(t, resp.Statuses, statuses)
}
