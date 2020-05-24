package server

import (
	"log"
	"net/http"

	"github.com/asankov/gira/pkg/client"
	"github.com/golangcollege/sessions"
)

type Page string

var (
	homePage       Page = "home.page.tmpl"
	listGamesPage  Page = "list.page.tmpl"
	createGamePage Page = "create.page.tmpl"
	signupUserPage Page = "signup.page.tmpl"
	loginUserPage  Page = "login.page.tmpl"
)

type Renderer interface {
	Render(w http.ResponseWriter, r *http.Request, d interface{}, p Page) error
}
type Server struct {
	Log      *log.Logger
	Session  *sessions.Session
	Client   *client.Client
	Renderer Renderer
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.routes().ServeHTTP(w, r)
}
