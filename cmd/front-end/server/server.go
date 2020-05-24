package server

import (
	"log"
	"net/http"

	"github.com/asankov/gira/pkg/client"
	"github.com/golangcollege/sessions"
)

var (
	homePage       = "home.page.tmpl"
	listGamesPage  = "list.page.tmpl"
	createGamePage = "create.page.tmpl"
	signupUserPage = "signup.page.tmpl"
	loginUserPage  = "login.page.tmpl"
)

// Renderer is the interface that will be used to interact with the part of the program
// that is responsible for rendering the web pages
type Renderer interface {
	Render(w http.ResponseWriter, r *http.Request, d interface{}, p string) error
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
