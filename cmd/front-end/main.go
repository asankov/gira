package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/asankov/gira/cmd/front-end/templates"
	"github.com/sirupsen/logrus"

	"github.com/asankov/gira/cmd/front-end/server"
	"github.com/asankov/gira/pkg/client"

	"github.com/golangcollege/sessions"
)

func main() {
	if err := run(); err != nil {
		logrus.Panic("error while running front-end service: " + err.Error())
	}
}

func run() error {
	port := flag.Int("port", 4001, "port on which the application is exposed")
	backEndAddr := flag.String("api_addr", "http://localhost:4000", "the address to the API service")
	sessionSecret := flag.String("session_secret", "s6Ndh+pPbnzHb7*297k1q5W0Tzbpa@ge", "32-byte secret that is to be used for the session store")
	logL := flag.String("log_level", "info", "the level of logging")
	flag.Parse()

	session := sessions.New([]byte(*sessionSecret))
	session.Lifetime = 12 * time.Hour

	cl, err := client.New(*backEndAddr)
	if err != nil {
		return fmt.Errorf("error while creating back-end client: %w", err)
	}

	log := logrus.New()
	logLevel, err := logrus.ParseLevel(*logL)
	if err != nil {
		return err
	}
	log.SetLevel(logLevel)
	logrus.SetLevel(logLevel)

	s := &server.Server{
		Log:      log,
		Client:   cl,
		Session:  session,
		Renderer: templates.NewRenderer(),
	}

	s.Log.Infoln(fmt.Sprintf("Front-end listening on port %d", *port))
	if err := http.ListenAndServe(fmt.Sprintf(":%d", *port), s); err != nil {
		return fmt.Errorf("error while listening: %w", err)
	}

	return nil
}
