package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/gira-games/client/pkg/client"

	"github.com/asankov/gira/cmd/front-end/templates"
	"github.com/sirupsen/logrus"

	"github.com/asankov/gira/cmd/front-end/server"

	"github.com/golangcollege/sessions"
)

func main() {
	if err := run(); err != nil {
		logrus.Panic("error while running front-end service: " + err.Error())
	}
}

func run() error {
	var (
		logL          = flag.String("log_level", "info", "the level of logging")
		port          = flag.Int("port", 4001, "port on which the application is exposed")
		backEndAddr   = flag.String("api_addr", "http://localhost:4000", "the address to the API service")
		sessionSecret = flag.String("session_secret", "s6Ndh+pPbnzHb7*297k1q5W0Tzbpa@ge", "32-byte secret that is to be used for the session store")
		enforceHTTPS  = flag.Bool("enforce-https", false, "whether or not to serve front-end via HTTPS")
	)
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

	addr := fmt.Sprintf(":%d", *port)
	srv := &http.Server{
		Addr:    addr,
		Handler: s,
		TLSConfig: &tls.Config{
			CurvePreferences:         []tls.CurveID{tls.CurveP256},
			PreferServerCipherSuites: true,
		},
	}

	if *enforceHTTPS {
		s.Log.Infoln(fmt.Sprintf("Front-end listening on %s via HTTPS", addr))
		err = srv.ListenAndServeTLS("tls/cert.pem", "tls/key.pem")
	} else {
		s.Log.Infoln(fmt.Sprintf("Front-end listening on %s", addr))
		err = srv.ListenAndServe()
	}

	if err != nil {
		return fmt.Errorf("error while listening: %w", err)
	}

	return nil
}
