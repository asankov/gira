package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/asankov/gira/pkg/models/postgres"

	// to register PostreSQL driver
	_ "github.com/lib/pq"
)

type server struct {
	log       *log.Logger
	handler   http.Handler
	gameModel *postgres.GameModel
	userModel *postgres.UserModel
}

func main() {
	if err := run(); err != nil {
		log.Panic("error while running server: " + err.Error())
	}
}

func run() error {
	port := *flag.Int("port", 4000, "port on which the application is exposed")
	dbHost := *flag.String("db_host", "localhost", "the address of the database")
	dbPort := *flag.Int("db_port", 5432, "the port of the database")
	dbUser := *flag.String("db_user", "antonsankov", "the user of the database")
	dbPass := *flag.String("db_pass", "", "the password for the database")
	dbName := *flag.String("db_name", "gira", "the name of the database")
	flag.Parse()

	db, err := openDB(dbHost, dbPort, dbUser, dbName, dbPass)
	if err != nil {
		return fmt.Errorf("error while opening DB: %w", err)
	}
	defer db.Close()

	s := &server{
		log:       log.New(os.Stdout, "", log.Ldate|log.Ltime),
		gameModel: &postgres.GameModel{DB: db},
		userModel: &postgres.UserModel{DB: db},
	}

	log.Println(fmt.Sprintf("listening on port %d", port))
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), s.routes()); err != nil {
		return fmt.Errorf("error while serving: %v", err)
	}

	return nil
}

func openDB(host string, port int, user string, dbName string, dbPass string) (*sql.DB, error) {
	connString := fmt.Sprintf("host=%s port=%d user=%s dbname=%s", host, port, user, dbName)
	if dbPass != "" {
		connString += fmt.Sprintf(" password=%s", dbPass)
	}
	connString += " sslmode=disable"
	fmt.Println("connString")
	fmt.Println(connString)
	db, err := sql.Open("postgres", connString)
	if err != nil {
		return nil, fmt.Errorf("error while opening connection to db: %w", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error while pinging db: %w", err)
	}

	return db, nil
}
