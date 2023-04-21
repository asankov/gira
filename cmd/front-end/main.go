package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"time"

	"github.com/asankov/gira/pkg/client"

	"github.com/asankov/gira/cmd/front-end/config"
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
	config, err := config.NewFromEnv()
	if err != nil {
		return fmt.Errorf("error while loading config: %w", err)
	}

	session := sessions.New([]byte(config.SessionSecret))
	session.Lifetime = 12 * time.Hour

	cl, err := client.New(config.APIAddress)
	if err != nil {
		return fmt.Errorf("error while creating back-end client: %w", err)
	}

	log := logrus.New()
	logLevel, err := logrus.ParseLevel(config.LogLevel)
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

	addr := fmt.Sprintf(":%d", config.Port)
	srv := &http.Server{
		Addr:    addr,
		Handler: s,
		TLSConfig: &tls.Config{
			CurvePreferences:         []tls.CurveID{tls.CurveP256},
			PreferServerCipherSuites: true,
		},
	}

	if config.EnforceHTTPS {
		s.Log.Infoln(fmt.Sprintf("Front-end listening on %s via HTTPS", addr))
		// TODO: these should be configurable
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
