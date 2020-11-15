package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
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

// Game is the struct that represents a game
type Game struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	FranchiseID string `json:"franchiseId"`
}

// GetGamesRequest is used when the consumer wants to get all games
type GetGamesRequest struct {
	Token           string
	ExcludeAssigned bool
}

// GetGamesResponse is the response that is returned from GetGames
type GetGamesResponse struct {
	Games []*Game
}

// CreateGameRequest is used when the consumer wants to create a games
type CreateGameRequest struct {
	Token string
	Game  *Game
}

// CreateGameResponse is the response that is returned from CreateGame
type CreateGameResponse struct {
	Game *Game
}

// GetGames returns all the games or all the games that are not assigned to the user
// to whom the token belongs.
func (c *Client) GetGames(ctx context.Context, request *GetGamesRequest) (*GetGamesResponse, error) {
	url := fmt.Sprintf("%s/games", c.addr)
	if request.ExcludeAssigned {
		url += "?excludeAssigned=true"
	}
	req, err := http.NewRequest(http.MethodGet, url, nil)
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

	var games GetGamesResponse
	if err := json.NewDecoder(res.Body).Decode(&games); err != nil {
		return nil, fmt.Errorf("error while decoding body: %w", err)
	}

	return &games, nil
}

// CreateGame creates a new game from the passed model.
func (c *Client) CreateGame(ctx context.Context, request *CreateGameRequest) (*CreateGameResponse, error) {
	body, err := json.Marshal(request.Game)
	if err != nil {
		return nil, ErrCreatingGame
	}
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/games", c.addr), bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("error while building HTTP request")
	}
	req.Header.Add(XAuthToken, request.Token)
	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, ErrCreatingGame
	}

	if res.StatusCode != http.StatusOK {
		if res.StatusCode == http.StatusUnauthorized {
			return nil, ErrNoAuthorization
		}
		// TODO
		// if res.StatusCode == http.StatusBadRequest {
		// 	// var jsonErr models.ErrorResponse
		// 	// if err := json.NewDecoder(res.Body).Decode(&jsonErr); err == nil {
		// 	// 	return nil, errors.New(jsonErr.Error)
		// 	// }
		// }
		return nil, ErrCreatingGame
	}

	var game Game
	if err := json.NewDecoder(res.Body).Decode(&game); err != nil {
		return nil, fmt.Errorf("error while decoding body: %w", err)
	}

	return &CreateGameResponse{Game: &game}, nil
}
