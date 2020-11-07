package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gira-games/api/pkg/models"
)

// ErrorResponse - this is duplicated with api/server/users_handlers.go
type ErrorResponse struct {
	Err string `json:"error"`
}

func (e *ErrorResponse) Error() string {
	return e.Err
}

func (c *Client) GetUser(token string) (*models.User, error) {
	url := fmt.Sprintf("%s/users", c.addr)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("error while building request")
	}
	req.Header.Set(models.XAuthToken, token)
	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error while calling %s: %w", url, err)
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error response from server: %d - %s", res.StatusCode, parseErrorBody(res))
	}
	var userResponse *models.UserResponse
	if err := json.NewDecoder(res.Body).Decode(&userResponse); err != nil {
		return nil, fmt.Errorf("error while decoding body: %w", err)
	}

	return userResponse.User, nil
}

func (c *Client) CreateUser(user *models.User) (*models.User, error) {
	body, err := json.Marshal(user)
	if err != nil {
		return nil, fmt.Errorf("error while building body: %w", err)
	}
	url := fmt.Sprintf("%s/users", c.addr)
	res, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, &ErrorResponse{Err: err.Error()}
	}
	if res.StatusCode != http.StatusOK {
		return nil, parseError(res)
		// return nil, fmt.Errorf("error response from server: %d - %s", res.StatusCode, parseErrorBody(res))
	}

	var userResponse *models.User
	if err := json.NewDecoder(res.Body).Decode(&userResponse); err != nil {
		return nil, fmt.Errorf("error while decoidng body: %w", err)
	}

	return userResponse, nil
}

func (c *Client) LoginUser(user *models.User) (*models.UserLoginResponse, error) {
	body, err := json.Marshal(user)
	if err != nil {
		return nil, fmt.Errorf("error while building body: %w", err)
	}
	url := fmt.Sprintf("%s/users/login", c.addr)
	res, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("error while calling %s: %w", url, err)
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error response from server: %d - %s", res.StatusCode, parseErrorBody(res))
	}

	var userResponse *models.UserLoginResponse
	if err := json.NewDecoder(res.Body).Decode(&userResponse); err != nil {
		return nil, fmt.Errorf("error while decoidng body: %w", err)
	}

	return userResponse, nil
}

func (c *Client) LogoutUser(token string) error {
	url := fmt.Sprintf("%s/users/logout", c.addr)
	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return fmt.Errorf("error while building request")
	}
	req.Header.Set(models.XAuthToken, token)
	res, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("error while calling %s: %w", url, err)
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("error response from server: %d - %s", res.StatusCode, parseErrorBody(res))
	}

	return nil
}

func parseError(r *http.Response) *ErrorResponse {
	var err ErrorResponse
	if err := json.NewDecoder(r.Body).Decode(&err); err != nil {
		return nil
	}
	return &err
}

func parseErrorBody(r *http.Response) string {
	errorBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return "<UNABLE TO READ ERROR BODY>"
	}
	return string(errorBody)
}
