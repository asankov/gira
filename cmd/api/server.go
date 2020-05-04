package main

import (
	"log"
	"net/http"

	"github.com/asankov/gira/internal/auth"
	"github.com/asankov/gira/pkg/models/postgres"
)

type server struct {
	log       *log.Logger
	gameModel *postgres.GameModel
	userModel *postgres.UserModel
	auth      *auth.Authenticator
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.routes().ServeHTTP(w, r)
}
