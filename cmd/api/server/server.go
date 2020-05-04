package server

import (
	"log"
	"net/http"

	"github.com/asankov/gira/internal/auth"
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

// Server is the struct that holds all the dependencies
// needed to run the application
type Server struct {
	Log       *log.Logger
	GameModel GameModel
	UserModel UserModel
	Auth      *auth.Authenticator
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.routes().ServeHTTP(w, r)
}
