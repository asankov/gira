package client

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/asankov/gira/pkg/models"
)

// GetUserGames returns all the games for the given user.
func (c *Client) GetUserGames(token string) ([]*models.Game, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/users/games", c.addr), nil)
	if err != nil {
		return nil, fmt.Errorf("error while building HTTP request")
	}
	req.Header.Add("x-auth-token", token)
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

	var games []*models.Game
	if err := json.NewDecoder(res.Body).Decode(&games); err != nil {
		return nil, fmt.Errorf("error while decoding body: %w", err)
	}

	return games, nil
}
