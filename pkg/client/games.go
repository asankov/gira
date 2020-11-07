package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/gira-games/api/pkg/models"
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
	// ErrChangingGame is returned when an error ocurred while changing status of game
	ErrChangingGame = errors.New("error while changing game")
	// ErrDeletingGame is returned when an error ocurred while deleting a game
	ErrDeletingGame = errors.New("error while deleting game")
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

type GetGamesOptions struct {
	ExcludeAssigned bool
}

var NoOptions = &GetGamesOptions{}

// GetGames returns all the games.
func (c *Client) GetGames(token string, options *GetGamesOptions) ([]*models.Game, error) {
	url := fmt.Sprintf("%s/games", c.addr)
	if options == nil {
		options = NoOptions
	}
	if options.ExcludeAssigned {
		url += "?excludeAssigned=true"
	}
	req, err := http.NewRequest(http.MethodGet, url, nil)
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

	var games models.GamesResponse
	if err := json.NewDecoder(res.Body).Decode(&games); err != nil {
		return nil, fmt.Errorf("error while decoding body: %w", err)
	}

	return games.Games, nil
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
	req.Header.Add(models.XAuthToken, token)
	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, ErrCreatingGame
	}

	if res.StatusCode != http.StatusOK {
		if res.StatusCode == http.StatusUnauthorized {
			return nil, ErrNoAuthorization
		}
		if res.StatusCode == http.StatusBadRequest {
			var jsonErr models.ErrorResponse
			if err := json.NewDecoder(res.Body).Decode(&jsonErr); err == nil {
				return nil, errors.New(jsonErr.Error)
			}
		}
		return nil, ErrCreatingGame
	}

	var gameResponse models.Game
	if err := json.NewDecoder(res.Body).Decode(&gameResponse); err != nil {
		return nil, fmt.Errorf("error while decoding body: %w", err)
	}

	return &gameResponse, nil
}
