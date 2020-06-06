package server

import (
	"log"
	"net/http"

	"github.com/asankov/gira/pkg/models"
)

// GameModel is the interface to interact with the Games provider (DB, service, etc.)
type GameModel interface {
	All() ([]*models.Game, error)
	Get(id string) (*models.Game, error)
	Insert(game *models.Game) (*models.Game, error)
}

// UserModel is the interface to interact with the User provider (DB, service, etc.)
type UserModel interface {
	Insert(user *models.User) (*models.User, error)
	Authenticate(email, password string) (*models.User, error)
}

// UserGamesModel is the interface to interact with the Users-Games relationship provider (DB, service, etc.)
type UserGamesModel interface {
	LinkGameToUser(userID, gameID string) error
	ChangeGameStatus(userID, userGameID string, status models.Status) error
	GetAvailableGamesFor(userID string) ([]*models.Game, error)
	GetUserGames(userID string) ([]*models.UserGame, error)
	GetUserGamesGrouped(userID string) (map[models.Status][]*models.UserGame, error)
}

// Authenticator is the interface to interact with the Authenticator (DB, OIDC provider, etc.)
type Authenticator interface {
	DecodeToken(token string) (*models.User, error)
	NewTokenForUser(user *models.User) (string, error)
}

// Server is the struct that holds all the dependencies
// needed to run the application
type Server struct {
	Log *log.Logger

	Authenticator
	GameModel
	UserModel
	UserGamesModel
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.routes().ServeHTTP(w, r)
}
