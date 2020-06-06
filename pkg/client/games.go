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
	// ErrNoAuthorization is returned when no authorization is sent for a authorized routes
	ErrNoAuthorization = errors.New("no authorization is present")
	// ErrLinkingGame is returned when an error occurred while linking game to user
	ErrLinkingGame = errors.New("error while linking game to user")
)

// Client is the struct that is used to communicate
// with the games service.
type Client struct {
	addr       string
	httpClient *http.Client
}

// New returns a new client with the given address.
func New(addr string) (*Client, error) {
	return &Client{
		addr:       addr,
		httpClient: &http.Client{},
	}, nil
}

// GetGames returns all the games.
func (c *Client) GetGames(token string) ([]*models.Game, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/games", c.addr), nil)
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

	var games models.GamesResponse
	if err := json.NewDecoder(res.Body).Decode(&games); err != nil {
		return nil, fmt.Errorf("error while decoidng body: %w", err)
	}

	return games.Games, nil
}

// GetGameByID returns the game with the given ID.
func (c *Client) GetGameByID(id string) (*models.Game, error) {
	return nil, nil

}

// CreateGame creates a new game from the passed model.
func (c *Client) CreateGame(game *models.Game, token string) (*models.Game, error) {
	body, err := json.Marshal(game)
	if err != nil {
		return nil, ErrCreatingGame
	}
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/games", c.addr), bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("error while building HTTP request")
	}
	req.Header.Add("x-auth-token", token)
	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, ErrCreatingGame
	}

	if res.StatusCode != http.StatusOK {
		if res.StatusCode == http.StatusUnauthorized {
			return nil, ErrNoAuthorization
		}
		return nil, ErrCreatingGame
	}

	var gameResponse models.Game
	if err := json.NewDecoder(res.Body).Decode(&gameResponse); err != nil {
		return nil, fmt.Errorf("error while decoidng body: %w", err)
	}

	return &gameResponse, nil
}
