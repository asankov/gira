package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golangcollege/sessions"
)

type server struct {
	log         *log.Logger
	backEndAddr string
	session     *sessions.Session
}

func main() {

	// TODO: replace this with configuration
	session := sessions.New([]byte("s6Ndh+pPbnzHb7*297k1q5W0Tzbpa@ge"))
	session.Lifetime = 12 * time.Hour

	s := &server{
		log: log.New(os.Stdout, "", log.Ldate|log.Ltime),
		// TODO: replace this with configuration
		// TODO: replace this with full-fledged client
		backEndAddr: "http://localhost:4000",
		session:     session,
	}

	s.log.Println("Front-end listening on port 4001")
	if err := http.ListenAndServe(":4001", s.routes()); err != nil {
		log.Fatalf("error while listening: %v", err)
	}
}
