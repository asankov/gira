package main

import (
	"fmt"

	"github.com/asankov/gira/cmd/api/config"
	"github.com/asankov/gira/cmd/api/database"

	"github.com/asankov/gira/cmd/api/server"
	"github.com/sirupsen/logrus"

	"github.com/asankov/gira/internal/auth"
	"github.com/asankov/gira/pkg/models/postgres"

	// to register PostreSQL driver
	_ "github.com/lib/pq"
)

func main() {
	if err := run(); err != nil {
		logrus.Fatalln("error while running server: " + err.Error())
	}
}

func run() error {
	config, err := config.NewFromEnv()
	if err != nil {
		return fmt.Errorf("error while loading config: %w", err)
	}

	log := logrus.New()
	logLevel, err := logrus.ParseLevel(config.LogLevel)
	if err != nil {
		return err
	}
	log.SetLevel(logLevel)
	logrus.SetLevel(logLevel)

	db, err := database.NewDB(&database.DBOptions{
		Host:   config.DB.Host,
		Port:   config.DB.Port,
		User:   config.DB.User,
		DBName: config.DB.Name,
		DBPass: config.DB.Password,
		UseSSL: config.UseSSL,
	})
	if err != nil {
		return fmt.Errorf("error while opening DB: %w", err)
	}
	defer db.Close()

	s := &server.Server{
		Log:            log,
		GameModel:      postgres.NewGameModel(db),
		UserModel:      postgres.NewUserModel(db),
		FranchiseModel: postgres.NewFranchiseModel(db),
		Authenticator:  auth.NewAutheniticator(config.Secret),
	}

	if err := s.Start(config.Port); err != nil {
		return fmt.Errorf("error while serving: %v", err)
	}

	return nil
}
