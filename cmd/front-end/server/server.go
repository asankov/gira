package server

import (
	"log"
	"net/http"

	"github.com/asankov/gira/pkg/client"
	"github.com/golangcollege/sessions"
)

type Server struct {
	Log     *log.Logger
	Session *sessions.Session
	Client  *client.Client
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.routes().ServeHTTP(w, r)
}
