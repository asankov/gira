package client

import (
	"github.com/asankov/gira/pkg/models"
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
	return nil, nil
}

// GetGameByID returns the game with the given ID.
func (c *Client) GetGameByID(id string) (*models.Game, error) {
	return nil, nil

}

// CreateGame creates a new game from the passed model.
func (c *Client) CreateGame(game models.Game) (*models.Game, error) {
	return nil, nil
}
