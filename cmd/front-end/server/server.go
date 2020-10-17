package server

import (
	"net/http"

	"github.com/asankov/gira/pkg/client"
	"github.com/asankov/gira/pkg/models"
	"github.com/golangcollege/sessions"
	"github.com/sirupsen/logrus"
)

var (
	homePage       = "home.page.tmpl"
	addGamePage    = "add.page.tmpl"
	listGamesPage  = "list.page.tmpl"
	createGamePage = "create.page.tmpl"
	signupUserPage = "signup.page.tmpl"
	loginUserPage  = "login.page.tmpl"

	emptyTemplateData = TemplateData{}
)

// TemplateData is the struct that holds all the data that can be passed to the template renderer to render
type TemplateData struct {
	Game       *models.Game
	User       *models.User
	Games      []*models.Game
	UserGames  []*models.UserGame
	Statuses   []models.Status
	Franchises []*models.Franchise

	Error string
	Flash string
}

// Renderer is the interface that will be used to interact with the part of the program
// that is responsible for rendering the web pages
type Renderer interface {
	Render(w http.ResponseWriter, r *http.Request, d TemplateData, p string) error
}

// APIClient is the interface that interacts with the API
type APIClient interface {
	LogoutUser(token string) error
	DeleteUserGame(gameID, token string) error
	GetUser(token string) (*models.User, error)
	CreateUser(user *models.User) (*models.User, error)
	LinkGameToUser(gameID, token string) (*models.UserGame, error)
	LoginUser(user *models.User) (*models.UserLoginResponse, error)
	CreateGame(game *models.Game, token string) (*models.Game, error)
	ChangeGameStatus(gameID, token string, status models.Status) error
	ChangeGameProgress(gameID, token string, progress *models.UserGameProgress) error
	GetUserGames(token string) (map[models.Status][]*models.UserGame, error)
	GetGames(token string, options *client.GetGamesOptions) ([]*models.Game, error)
	GetFranchises(token string) ([]*models.Franchise, error)
}

// Server is the struct that holds all the dependencies
// needed to run the application
type Server struct {
	Log      *logrus.Logger
	Session  *sessions.Session
	Client   APIClient
	Renderer Renderer
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.routes().ServeHTTP(w, r)
}
