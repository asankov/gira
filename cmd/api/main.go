package main

import (
	"flag"
	"fmt"

	"github.com/gira-games/api/cmd/api/database"

	"github.com/gira-games/api/cmd/api/server"
	"github.com/sirupsen/logrus"

	"github.com/gira-games/api/internal/auth"
	"github.com/gira-games/api/pkg/models/postgres"

	// to register PostreSQL driver
	_ "github.com/lib/pq"
)

func main() {
	if err := run(); err != nil {
		logrus.Fatalln("error while running server: " + err.Error())
	}
}

func run() error {
	port := flag.Int("port", 4000, "port on which the application is exposed")
	dbHost := flag.String("db_host", "localhost", "the address of the database")
	dbPort := flag.Int("db_port", 5432, "the port of the database")
	dbUser := flag.String("db_user", "antonsankov", "the user of the database")
	dbPass := flag.String("db_pass", "", "the password for the database")
	dbName := flag.String("db_name", "gira", "the name of the database")
	secret := flag.String("token_string", "9^ahslgndb&ahas2ey*hasdh732rbusd", "secret to be used for encoding and decoding JWT tokens")
	useSSL := flag.Bool("use_ssl", false, "whether or not to use SSL when connecting to DB")
	logL := flag.String("log_level", "info", "the level of logging")
	flag.Parse()

	log := logrus.New()
	logLevel, err := logrus.ParseLevel(*logL)
	if err != nil {
		return err
	}
	log.SetLevel(logLevel)
	logrus.SetLevel(logLevel)

	db, err := database.NewDB(&database.DBOptions{
		Host:   *dbHost,
		Port:   *dbPort,
		User:   *dbUser,
		DBName: *dbName,
		DBPass: *dbPass,
		UseSSL: *useSSL,
	})
	if err != nil {
		return fmt.Errorf("error while opening DB: %w", err)
	}
	defer db.Close()

	s := &server.Server{
		Log:            log,
		GameModel:      &postgres.GameModel{DB: db},
		UserModel:      &postgres.UserModel{DB: db},
		UserGamesModel: &postgres.UserGamesModel{DB: db},
		FranchiseModel: &postgres.FranchiseModel{DB: db},
		Authenticator:  auth.NewAutheniticator(*secret),
	}

	if err := s.Start(*port); err != nil {
		return fmt.Errorf("error while serving: %v", err)
	}

	return s.Start(*port)
}
