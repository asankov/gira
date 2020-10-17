package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/asankov/gira/pkg/models"
)

var (
	// ErrFetchingFranchises is a generic error
	ErrFetchingFranchises = errors.New("error while fetching franchises")
)

func (c *Client) GetFranchises(token string) ([]*models.Franchise, error) {
	url := fmt.Sprintf("%s/franchises", c.addr)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("error while building HTTP request")
	}
	req.Header.Add(models.XAuthToken, token)
	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, ErrFetchingFranchises
	}
	if res.StatusCode != http.StatusOK {
		if res.StatusCode == http.StatusUnauthorized {
			return nil, ErrNoAuthorization
		}
		return nil, ErrFetchingFranchises
	}

	var franchises models.FranchisesResponse
	if err := json.NewDecoder(res.Body).Decode(&franchises); err != nil {
		return nil, fmt.Errorf("error while decoding body: %w", err)
	}

	return franchises.Franchises, nil
}
