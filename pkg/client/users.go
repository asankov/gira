package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/asankov/gira/pkg/models"
)

func (c *Client) CreateUser(user *models.User) (*models.User, error) {
	body, err := json.Marshal(user)
	if err != nil {
		return nil, fmt.Errorf("error while building body: %w", err)
	}
	url := fmt.Sprintf("%s/users", c.addr)
	res, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("error while calling %s: %w", url, err)
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error response from server: %d - %s", res.StatusCode, parseErrorBody(res))
	}

	var userResponse *models.User
	json.NewDecoder(res.Body).Decode(&userResponse)

	return userResponse, nil
}

func (c *Client) LoginUser(user *models.User) (*models.UserResponse, error) {
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

	var userResponse *models.UserResponse
	json.NewDecoder(res.Body).Decode(&userResponse)

	return userResponse, nil
}

func parseErrorBody(r *http.Response) string {
	errorBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return "<UNABLE TO READ ERROR BODY>"
	}
	return string(errorBody)
}
