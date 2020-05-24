package server

import (
	"log"
	"net/http"

	"github.com/asankov/gira/pkg/client"
	"github.com/golangcollege/sessions"
)

// Page represents a Page to be rendered
type Page string

var (
	homePage       Page = "home.page.tmpl"
	listGamesPage  Page = "list.page.tmpl"
	createGamePage Page = "create.page.tmpl"
	signupUserPage Page = "signup.page.tmpl"
	loginUserPage  Page = "login.page.tmpl"
)

// Renderer is the interface that will be used to interact with the part of the program
// that is responsible for rendering the web pages
type Renderer interface {
	Render(w http.ResponseWriter, r *http.Request, d interface{}, p Page) error
}

// Server is the struct that holds all the dependencies
// needed to run the application
type Server struct {
	Log      *log.Logger
	Session  *sessions.Session
	Client   *client.Client
	Renderer Renderer
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.routes().ServeHTTP(w, r)
}
