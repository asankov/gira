package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// User is the representation of a user
type User struct {
	ID             string `json:"id,omitempty"`
	Username       string `json:"username,omitempty"`
	Email          string `json:"email,omitempty"`
	Password       string `json:"password,omitempty"`
	HashedPassword []byte `json:"-"`
}

type GetUserRequest struct {
	Token string
}

type GetUserResponse struct {
	ID       string `json:"id,omitempty"`
	Username string `json:"username,omitempty"`
	Email    string `json:"email,omitempty"`
}

type CreateUserRequest struct {
	Username string `json:"username,omitempty"`
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

type CreateUserResponse struct {
	ID       string `json:"id,omitempty"`
	Username string `json:"username,omitempty"`
	Email    string `json:"email,omitempty"`
}

type LoginUserRequest struct {
	Username string `json:"username,omitempty"`
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

// UserLoginResponse is the response that is returned
// when a user is logged in.
type UserLoginResponse struct {
	Token string `json:"token"`
}

type LogoutUserRequest struct {
	Token string
}

// ErrorResponse - this is duplicated with api/server/users_handlers.go
type ErrorResponse struct {
	Err string `json:"error"`
}

func (e *ErrorResponse) Error() string {
	return e.Err
}

func (c *Client) GetUser(ctx context.Context, request *GetUserRequest) (*GetUserResponse, error) {
	url := fmt.Sprintf("%s/users", c.addr)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("error while building request")
	}
	req.Header.Set(XAuthToken, request.Token)
	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error while calling %s: %w", url, err)
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error response from server: %d - %s", res.StatusCode, parseErrorBody(res))
	}
	var userResponse struct {
		User *GetUserResponse `json:"user"`
	}
	if err := json.NewDecoder(res.Body).Decode(&userResponse); err != nil {
		return nil, fmt.Errorf("error while decoding body: %w", err)
	}

	return userResponse.User, nil
}

func (c *Client) CreateUser(ctx context.Context, request *CreateUserRequest) (*CreateUserResponse, error) {
	body, err := json.Marshal(request)
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

	var userResponse *CreateUserResponse
	if err := json.NewDecoder(res.Body).Decode(&userResponse); err != nil {
		return nil, fmt.Errorf("error while decoidng body: %w", err)
	}

	return userResponse, nil
}

func (c *Client) LoginUser(ctx context.Context, request *LoginUserRequest) (*UserLoginResponse, error) {
	body, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("error while building body: %w", err)
	}
	url := fmt.Sprintf("%s/users/login", c.addr)
	res, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("error while calling %s: %w", url, err)
	}
	if res.StatusCode != http.StatusOK {
		return nil, &ErrorResponse{Err: parseErrorBody(res)}
	}

	var userResponse *UserLoginResponse
	if err := json.NewDecoder(res.Body).Decode(&userResponse); err != nil {
		return nil, fmt.Errorf("error while decoidng body: %w", err)
	}

	return userResponse, nil
}

func (c *Client) LogoutUser(ctx context.Context, request *LogoutUserRequest) error {
	url := fmt.Sprintf("%s/users/logout", c.addr)
	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return fmt.Errorf("error while building request")
	}
	req.Header.Set(XAuthToken, request.Token)
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
