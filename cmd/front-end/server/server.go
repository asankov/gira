package server

import (
	"context"
	"net/http"

	"github.com/asankov/gira/pkg/client"

	"github.com/golangcollege/sessions"
	"github.com/sirupsen/logrus"
)

var (
	homePage       = "home.page.tmpl"
	listGamesPage  = "list.page.tmpl"
	createGamePage = "create.page.tmpl"
	signupUserPage = "signup.page.tmpl"
	loginUserPage  = "login.page.tmpl"

	emptyTemplateData = TemplateData{}
)

// TemplateData is the struct that holds all the data that can be passed to the template renderer to render
type TemplateData struct {
	Game       *client.Game
	User       *client.User
	Games      []TemplateGame
	Statuses   []client.Status
	Franchises []*client.Franchise

	SelectedFranchiseID string
	Error               string
	Flash               string
}

// TemplateGame is the struct that holds all the game info that is passed to the template renderer to render
type TemplateGame struct {
	ID            string
	Name          string
	FranchiseID   string
	FranchiseName string

	Status   client.Status
	Progress *client.GameProgress
}

// Renderer is the interface that will be used to interact with the part of the program
// that is responsible for rendering the web pages
type Renderer interface {
	Render(w http.ResponseWriter, r *http.Request, d TemplateData, p string) error
}

// APIClient is the interface that interacts with the API
type APIClient interface {
	GetFranchises(context.Context, *client.GetFranchisesRequest) (*client.GetFranchisesResponse, error)
	CreateFranchise(context.Context, *client.CreateFranchiseRequest) (*client.CreateFranchiseResponse, error)

	GetGames(context.Context, *client.GetGamesRequest) (*client.GetGamesResponse, error)
	CreateGame(context.Context, *client.CreateGameRequest) (*client.CreateGameResponse, error)

	UpdateGameProgress(context.Context, *client.UpdateGameProgressRequest) error
	DeleteUserGame(context.Context, *client.DeleteUserGameRequest) error

	LoginUser(context.Context, *client.LoginUserRequest) (*client.UserLoginResponse, error)
	CreateUser(context.Context, *client.CreateUserRequest) (*client.CreateUserResponse, error)
	GetUser(context.Context, *client.GetUserRequest) (*client.GetUserResponse, error)
	LogoutUser(context.Context, *client.LogoutUserRequest) error

	GetStatuses(ctx context.Context, request *client.GetStatusesRequest) (*client.GetStatusesResponse, error)
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
