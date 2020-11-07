package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gira-games/api/pkg/models"
)

// GetUserGames returns all the games for the given user.
func (c *Client) GetUserGames(token string) (map[models.Status][]*models.UserGame, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/users/games", c.addr), nil)
	if err != nil {
		return nil, fmt.Errorf("error while building HTTP request")
	}
	req.Header.Add(models.XAuthToken, token)
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

	var gamesResponse map[models.Status][]*models.UserGame
	if err := json.NewDecoder(res.Body).Decode(&gamesResponse); err != nil {
		return nil, fmt.Errorf("error while decoding body: %w", err)
	}

	return gamesResponse, nil
}

// LinkGameToUser adds the game with the given ID to the list of games of the user, associated with the token
func (c *Client) LinkGameToUser(gameID, token string) (*models.UserGame, error) {
	body, err := json.Marshal(models.UserGameRequest{Game: &models.Game{ID: gameID}})
	if err != nil {
		return nil, fmt.Errorf("error while marshalling body: %w", err)
	}
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/users/games", c.addr), bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("error while building HTTP request")
	}
	req.Header.Add(models.XAuthToken, token)
	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, ErrLinkingGame
	}
	if res.StatusCode != http.StatusOK {
		if res.StatusCode == http.StatusUnauthorized {
			return nil, ErrNoAuthorization
		}
		return nil, ErrLinkingGame
	}

	// TODO: return real response
	return nil, nil
}

func (c *Client) ChangeGameStatus(gameID, token string, status models.Status) error {
	return c.changeGame(gameID, token, &models.ChangeGameStatusRequest{Status: status})

}

func (c *Client) ChangeGameProgress(gameID, token string, progress *models.UserGameProgress) error {
	return c.changeGame(gameID, token, &models.ChangeGameStatusRequest{Progress: progress})
}

func (c *Client) changeGame(gameID, token string, changeGameReq *models.ChangeGameStatusRequest) error {
	body, err := json.Marshal(changeGameReq)
	if err != nil {
		return fmt.Errorf("error while marshalling body: %w", err)
	}

	req, err := http.NewRequest(http.MethodPatch, fmt.Sprintf("%s/users/games/%s", c.addr, gameID), bytes.NewBuffer((body)))
	if err != nil {
		return fmt.Errorf("error while building HTTP request")
	}
	req.Header.Add(models.XAuthToken, token)
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

	// TODO: return real response
	return nil
}

func (c *Client) DeleteUserGame(gameID, token string) error {
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/users/games/%s", c.addr, gameID), nil)
	if err != nil {
		return fmt.Errorf("error while building HTTP request")
	}
	req.Header.Add(models.XAuthToken, token)
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
