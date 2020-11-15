package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

var (
	// ErrFetchingFranchises is a generic error
	ErrFetchingFranchises = errors.New("error while fetching franchises")
	// ErrCreatingFranchise is a generic error
	ErrCreatingFranchise = errors.New("error while creating franchise")
)

// Franchise is the struct that represents a franchise
type Franchise struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

// CreateFranchiseRequest is used when the consumer wants to create a franchise
type CreateFranchiseRequest struct {
	Name  string
	Token string
}

// CreateFranchiseResponse is the response that is returned from CreateFranchise
type CreateFranchiseResponse struct {
	Franchise *Franchise
}

// GetFranchisesRequest is used when the consumer wants to get all franchises
type GetFranchisesRequest struct {
	Token string
}

// GetFranchisesResponse is the response that is returned from GetFranchises
type GetFranchisesResponse struct {
	Franchises []*Franchise `json:"franchises,omitempty"`
}

// GetFranchises returns all the franchises
func (c *Client) GetFranchises(ctx context.Context, request *GetFranchisesRequest) (*GetFranchisesResponse, error) {
	url := fmt.Sprintf("%s/franchises", c.addr)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("error while building HTTP request")
	}
	req.Header.Add(XAuthToken, request.Token)
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

	var franchisesResponse GetFranchisesResponse
	if err := json.NewDecoder(res.Body).Decode(&franchisesResponse); err != nil {
		return nil, fmt.Errorf("error while decoding body: %w", err)
	}

	return &franchisesResponse, nil
}

// CreateFranchise creates a franchise
func (c *Client) CreateFranchise(ctx context.Context, req *CreateFranchiseRequest) (*CreateFranchiseResponse, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, ErrCreatingFranchise
	}
	url := fmt.Sprintf("%s/franchises", c.addr)
	request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("error while building HTTP request")
	}
	request.Header.Add(XAuthToken, req.Token)
	res, err := c.httpClient.Do(request)
	if err != nil {
		return nil, ErrCreatingFranchise
	}
	if res.StatusCode != http.StatusOK {
		if res.StatusCode == http.StatusUnauthorized {
			return nil, ErrNoAuthorization
		}
		// TODO:

		// if res.StatusCode == http.StatusBadRequest {
		// 	var jsonErr models.ErrorResponse
		// 	if err := json.NewDecoder(res.Body).Decode(&jsonErr); err != nil {
		// 		return nil, ErrCreatingFranchise
		// 	}
		// 	return nil, errors.New(jsonErr.Error)
		// }
		// return nil, ErrCreatingFranchise
	}

	var franchise Franchise
	if err := json.NewDecoder(res.Body).Decode(&franchise); err != nil {
		return nil, fmt.Errorf("error while decoding body: %w", err)
	}

	return &CreateFranchiseResponse{Franchise: &franchise}, nil
}
