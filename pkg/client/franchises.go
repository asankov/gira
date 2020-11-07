package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/gira-games/api/pkg/models"
)

var (
	// ErrFetchingFranchises is a generic error
	ErrFetchingFranchises = errors.New("error while fetching franchises")
	// ErrCreatingFranchise is a generic error
	ErrCreatingFranchise = errors.New("error while creating franchise")
)

type CreateFranchiseRequest struct {
	Name string `json:"name"`
}

type CreateFranchiseResponse struct {
	Franchise *models.Franchise
}

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

func (c *Client) CreateFranchise(req *CreateFranchiseRequest, token string) (*CreateFranchiseResponse, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, ErrCreatingFranchise
	}
	url := fmt.Sprintf("%s/franchises", c.addr)
	request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("error while building HTTP request")
	}
	request.Header.Add(models.XAuthToken, token)
	res, err := c.httpClient.Do(request)
	if err != nil {
		return nil, ErrCreatingFranchise
	}
	if res.StatusCode != http.StatusOK {
		if res.StatusCode == http.StatusUnauthorized {
			return nil, ErrNoAuthorization
		}
		if res.StatusCode == http.StatusBadRequest {
			var jsonErr models.ErrorResponse
			if err := json.NewDecoder(res.Body).Decode(&jsonErr); err != nil {
				return nil, ErrCreatingFranchise
			}
			return nil, errors.New(jsonErr.Error)
		}
		return nil, ErrCreatingFranchise
	}

	var franchise models.Franchise
	if err := json.NewDecoder(res.Body).Decode(&franchise); err != nil {
		return nil, fmt.Errorf("error while decoding body: %w", err)
	}

	return &CreateFranchiseResponse{Franchise: &franchise}, nil
}
