package main

import (
	"flag"
	"fmt"
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
	port := *flag.Int("port", 4001, "port on which the application is exposed")
	backEndAddr := *flag.String("api_addr", "http://localhost:4000", "the address to the API service")
	sessionSecret := *flag.String("session_secret", "s6Ndh+pPbnzHb7*297k1q5W0Tzbpa@ge", "32-byte secret that is to be used for the session store")
	flag.Parse()

	session := sessions.New([]byte(sessionSecret))
	session.Lifetime = 12 * time.Hour

	s := &server{
		log: log.New(os.Stdout, "", log.Ldate|log.Ltime),
		// TODO: replace this with full-fledged client
		backEndAddr: backEndAddr,
		session:     session,
	}

	s.log.Println(fmt.Sprintf("Front-end listening on port %d", port))
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), s.routes()); err != nil {
		log.Fatalf("error while listening: %v", err)
	}
}
