package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// UserGame is the representation of a user game relation
type UserGame struct {
	ID       string            `json:"id,omitempty"`
	User     *User             `json:"user,omitempty"`
	Game     *Game             `json:"game,omitempty"`
	Status   Status            `json:"status,omitempty"`
	Progress *UserGameProgress `json:"progress,omitempty"`
}

type UserGameProgress struct {
	Current int `json:"current,omitempty"`
	Final   int `json:"final,omitempty"`
}

type GetUserGamesRequest struct {
	Token string
}

type GetUserGamesResponse struct {
	UserGames map[Status][]*UserGame
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
	Status   Status            `json:"status,omitempty"`
	Progress *UserGameProgress `json:"progress,omitempty"`
}

type DeleteUserGameRequest struct {
	GameID string
	Token  string
}

// GetUserGames returns all the games for the given user.
func (c *Client) GetUserGames(ctx context.Context, request *GetUserGamesRequest) (*GetUserGamesResponse, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/users/games", c.addr), nil)
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
		return nil, ErrFetchingGames
	}

	var gamesResponse map[Status][]*UserGame
	if err := json.NewDecoder(res.Body).Decode(&gamesResponse); err != nil {
		return nil, fmt.Errorf("error while decoding body: %w", err)
	}

	return &GetUserGamesResponse{UserGames: gamesResponse}, nil
}

// LinkGameToUser adds the game with the given ID to the list of games of the user, associated with the token
func (c *Client) LinkGameToUser(ctx context.Context, request *LinkGameToUserRequest) error {
	body, err := json.Marshal(struct {
		Game struct {
			ID string `json:"id"`
		}
	}{
		Game: struct {
			ID string `json:"id"`
		}{
			ID: request.GameID,
		},
	})
	if err != nil {
		return fmt.Errorf("error while marshalling body: %w", err)
	}
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/users/games", c.addr), bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("error while building HTTP request")
	}
	req.Header.Add(XAuthToken, request.Token)
	res, err := c.httpClient.Do(req)
	if err != nil {
		return ErrLinkingGame
	}
	if res.StatusCode != http.StatusOK {
		if res.StatusCode == http.StatusUnauthorized {
			return ErrNoAuthorization
		}
		return ErrLinkingGame
	}

	return nil
}

func (c *Client) UpdateGameProgress(ctx context.Context, request *UpdateGameProgressRequest) error {
	body, err := json.Marshal(request.Update)
	if err != nil {
		return fmt.Errorf("error while marshalling body: %w", err)
	}

	req, err := http.NewRequest(http.MethodPatch, fmt.Sprintf("%s/users/games/%s", c.addr, request.GameID), bytes.NewBuffer((body)))
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
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/users/games/%s", c.addr, request.GameID), nil)
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
