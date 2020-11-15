package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

var (
	// ErrFetchingStatuses is a generic error
	ErrFetchingStatuses = errors.New("error while fetching statuses")
)

// Status is the status of the game
type Status string

// GetStatusesRequest is used when getting the statuses from the server
type GetStatusesRequest struct {
	Token string
}

// GetStatusesResponse is returned from the GetStatuses method
type GetStatusesResponse struct {
	Statuses []Status `json:"statuses,omitempty"`
}

// GetStatuses fetches the statuses from the server
func (c *Client) GetStatuses(ctx context.Context, request *GetStatusesRequest) (*GetStatusesResponse, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/statuses", c.addr), nil)
	if err != nil {
		return nil, fmt.Errorf("error while building HTTP request")
	}
	req.Header.Add(XAuthToken, request.Token)
	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, ErrFetchingGames
	}
	if res.StatusCode != http.StatusOK {
		if res.StatusCode == http.StatusUnauthorized {
			return nil, ErrNoAuthorization
		}
		return nil, ErrFetchingStatuses
	}

	var statusesResponse GetStatusesResponse
	if err := json.NewDecoder(res.Body).Decode(&statusesResponse); err != nil {
		return nil, fmt.Errorf("error while decoding body: %w", err)
	}
	return &statusesResponse, nil
}
