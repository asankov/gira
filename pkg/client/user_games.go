package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type GameProgress struct {
	Current int `json:"current,omitempty"`
	Final   int `json:"final,omitempty"`
}

type LinkGameToUserRequest struct {
	Token  string
	GameID string
}

type UpdateGameProgressRequest struct {
	GameID string
	Token  string
	Update UpdateGameProgressChange
}

type UpdateGameProgressChange struct {
	Status   Status        `json:"status,omitempty"`
	Progress *GameProgress `json:"progress,omitempty"`
}

type DeleteUserGameRequest struct {
	GameID string
	Token  string
}

func (c *Client) UpdateGameProgress(ctx context.Context, request *UpdateGameProgressRequest) error {
	body, err := json.Marshal(request.Update)
	if err != nil {
		return fmt.Errorf("error while marshalling body: %w", err)
	}

	req, err := http.NewRequest(http.MethodPatch, fmt.Sprintf("%s/games/%s", c.addr, request.GameID), bytes.NewBuffer((body)))
	if err != nil {
		return fmt.Errorf("error while building HTTP request")
	}
	req.Header.Add(XAuthToken, request.Token)
	res, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		if res.StatusCode == http.StatusUnauthorized {
			return ErrNoAuthorization
		}
		return ErrChangingGame
	}

	return nil
}

func (c *Client) DeleteUserGame(ctx context.Context, request *DeleteUserGameRequest) error {
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/games/%s", c.addr, request.GameID), nil)
	if err != nil {
		return fmt.Errorf("error while building HTTP request")
	}
	req.Header.Add(XAuthToken, request.Token)
	res, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		if res.StatusCode == http.StatusUnauthorized {
			return ErrNoAuthorization
		}
		return ErrDeletingGame
	}

	// TODO: return real response
	return nil
}
