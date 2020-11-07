package server

import (
	"fmt"
	"net/http"

	"github.com/gira-games/api/pkg/models"
	"github.com/sirupsen/logrus"
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
	AssociateTokenWithUser(userID, token string) error
	InvalidateToken(userID, token string) error
	GetUserByToken(token string) (*models.User, error)
}

// UserGamesModel is the interface to interact with the Users-Games relationship provider (DB, service, etc.)
type UserGamesModel interface {
	LinkGameToUser(userID, gameID string, progress *models.UserGameProgress) error
	ChangeGameStatus(userID, userGameID string, status models.Status) error
	ChangeGameProgress(userID, userGameID string, progress *models.UserGameProgress) error
	GetAvailableGamesFor(userID string) ([]*models.Game, error)
	GetUserGames(userID string) ([]*models.UserGame, error)
	GetUserGamesGrouped(userID string) (map[models.Status][]*models.UserGame, error)
	DeleteUserGame(userGameID string) error
}

// FranchiseModel is the interface to interact with the Franchise provider (DB, service, etc.)
type FranchiseModel interface {
	Insert(franchise *models.Franchise) (*models.Franchise, error)
	All() ([]*models.Franchise, error)
}

// Authenticator is the interface to interact with the Authenticator (DB, OIDC provider, etc.)
type Authenticator interface {
	DecodeToken(token string) (*models.User, error)
	NewTokenForUser(user *models.User) (string, error)
}

// Server is the struct that holds all the dependencies
// needed to run the application
type Server struct {
	Log *logrus.Logger

	Authenticator
	GameModel
	UserModel
	UserGamesModel
	FranchiseModel
}

// Options is the struct used to construct a server
type Options struct {
	Log *logrus.Logger

	Authenticator
	GameModel
	UserModel
	UserGamesModel
	FranchiseModel
}

// New returns a new Server, based on opts.
// It also validates the arguments and can return an error
// if it receives options that are not right.
func New(opts *Options) (*Server, error) {
	// TODO: validate args
	return &Server{
		Log:            opts.Log,
		Authenticator:  opts.Authenticator,
		GameModel:      opts.GameModel,
		UserModel:      opts.UserModel,
		UserGamesModel: opts.UserGamesModel,
		FranchiseModel: opts.FranchiseModel,
	}, nil
}

// ServeHTTP implement the http.Handler inteface
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.routes().ServeHTTP(w, r)
}

// Start starts the server listenning on the given port
func (s *Server) Start(port int) error {
	s.Log.Infoln(fmt.Sprintf("listening on port %d", port))
	return http.ListenAndServe(fmt.Sprintf(":%d", port), s)
}
