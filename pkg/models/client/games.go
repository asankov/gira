package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/asankov/gira/pkg/models"
)

var (
	// ErrFetchingGames is a generic error
	ErrFetchingGames = errors.New("error while fetching games")
)

// Client is the struct that is used to communicate
// with the games service.
type Client struct {
	addr string
}

// New returns a new client with the given address.
func New(addr string) (*Client, error) {
	return &Client{
		addr: addr,
	}, nil
}

// GetGames returns all the games.
func (c *Client) GetGames() ([]*models.Game, error) {
	res, err := http.Get(fmt.Sprintf("%s/games", c.addr))
	if err != nil {
		return nil, ErrFetchingGames
	}
	if res.StatusCode != 200 {
		return nil, ErrFetchingGames
	}

	var games []*models.Game
	json.NewDecoder(res.Body).Decode(&games)

	return games, nil
}

// GetGameByID returns the game with the given ID.
func (c *Client) GetGameByID(id string) (*models.Game, error) {
	return nil, nil

}

// CreateGame creates a new game from the passed model.
func (c *Client) CreateGame(game models.Game) (*models.Game, error) {
	return nil, nil
}
