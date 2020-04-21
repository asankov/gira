package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/asankov/gira/pkg/models"
)

var (
	// ErrFetchingGames is a generic error
	ErrFetchingGames = errors.New("error while fetching games")
	// ErrCreatingGame is a generic error
	ErrCreatingGame = errors.New("error while creating game")
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
func (c *Client) CreateGame(game *models.Game) (*models.Game, error) {
	body, err := json.Marshal(game)
	if err != nil {
		return nil, ErrCreatingGame
	}
	res, err := http.Post(fmt.Sprintf("%s/games", c.addr), "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, ErrCreatingGame
	}

	if res.StatusCode != http.StatusOK {
		return nil, ErrCreatingGame
	}

	var gameResponse models.Game
	json.NewDecoder(res.Body).Decode(&gameResponse)

	return &gameResponse, nil
}
