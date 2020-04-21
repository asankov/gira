package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/asankov/gira/pkg/models/client"

	"github.com/golangcollege/sessions"
)

type server struct {
	log         *log.Logger
	backEndAddr string
	session     *sessions.Session
	client      *client.Client
}

func main() {
	port := *flag.Int("port", 4001, "port on which the application is exposed")
	backEndAddr := *flag.String("api_addr", "http://localhost:4000", "the address to the API service")
	sessionSecret := *flag.String("session_secret", "s6Ndh+pPbnzHb7*297k1q5W0Tzbpa@ge", "32-byte secret that is to be used for the session store")
	flag.Parse()

	session := sessions.New([]byte(sessionSecret))
	session.Lifetime = 12 * time.Hour

	cl, err := client.New(backEndAddr)
	if err != nil {
		log.Fatalf("error while creating back-end client: %v", err)
	}

	s := &server{
		log:     log.New(os.Stdout, "", log.Ldate|log.Ltime),
		client:  cl,
		session: session,
	}

	s.log.Println(fmt.Sprintf("Front-end listening on port %d", port))
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), s.routes()); err != nil {
		log.Fatalf("error while listening: %v", err)
	}
}
