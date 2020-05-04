package main

import (
	"log"
	"net/http"

	"github.com/asankov/gira/pkg/client"
	"github.com/golangcollege/sessions"
)

type server struct {
	log     *log.Logger
	session *sessions.Session
	client  *client.Client
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.routes().ServeHTTP(w, r)
}
