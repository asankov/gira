package models

import (
	"fmt"
)

// Game is the representation of a game
// in the database.
type Game struct {
	ID          string        `json:"id"`
	Name        string        `json:"name,omitempty"`
	Franchise   string        `json:"franchise,omitempty"`
	FranchiseID string        `json:"franchiseId,omitempty"`
	Status      Status        `json:"status,omitempty"`
	Progress    *GameProgress `json:"progress,omitempty"`

	UserID string `json:"-"`
}

type GameProgress struct {
	Current int `json:"current"`
	Final   int `json:"final,omitempty"`
}

type GamesResponse struct {
	Games []*Game `json:"games"`
}

// User is the representation of a user
// in the database.
type User struct {
	ID             string `json:"id,omitempty"`
	Username       string `json:"username,omitempty"`
	Email          string `json:"email,omitempty"`
	Password       string `json:"password,omitempty"`
	HashedPassword []byte `json:"-"`
}

// Status is the type that represents the status of a game
type Status string

var (
	// StatusTODO is the To Do status of a game
	StatusTODO Status = "To Do"
	// StatusInProgress is the In progress status of a game
	StatusInProgress Status = "In Progress"
	// StatusDone is the Done status of the game
	StatusDone Status = "Done"

	// AllStatuses is collection of all statuses
	AllStatuses = []Status{
		StatusTODO,
		StatusInProgress,
		StatusDone,
	}
)

// Validate shows whether the status is a valid status
// and returns an error if not.
func (s Status) Validate() error {
	for _, status := range AllStatuses {
		if s == status {
			return nil
		}
	}
	return fmt.Errorf("%s is not a valid status", s)
}

// StatusesResponse is the response that is returned from the Statuses API
type StatusesResponse struct {
	Statuses []Status `json:"statuses,omitempty"`
}

// UserLoginResponse is the response that is returned
// when a user is logged in.
type UserLoginResponse struct {
	Token string `json:"token"`
}

// UserResponse is the response that is returned
// from the GET /users API
type UserResponse struct {
	User *User `json:"user"`
}

type UserGameRequest struct {
	Game     *Game         `json:"game"`
	Progress *GameProgress `json:"progress,omitempty"`
}

type UserGamesResponse struct {
	Games []*Game `json:"games"`
}

type UserGameResponse struct {
	ID     string `json:"id"`
	Game   *Game  `json:"game"`
	Status Status `json:"status"`
}

// ErrorResponse is the generic error response returned from the API,
// when an error of any kind occurred.
type ErrorResponse struct {
	Error string `json:"error,omitempty"`
}

type ChangeGameStatusRequest struct {
	Status   Status        `json:"status,omitempty"`
	Progress *GameProgress `json:"progress,omitempty"`
}

type Franchise struct {
	ID   string `json:"id"`
	Name string `json:"name"`

	UserID string `json:"-"`
}

type FranchisesResponse struct {
	Franchises []*Franchise `json:"franchises"`
}
