package main

import (
	"database/sql"
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
}

func main() {
	if err := run(); err != nil {
		log.Panic("error while running server: " + err.Error())
	}
}

func run() error {

	db, err := openDB("localhost", 5432, "antonsankov", "gira")
	if err != nil {
		return fmt.Errorf("error while opening DB: %w", err)
	}
	defer db.Close()

	s := &server{
		log:       log.New(os.Stdout, "", log.Ldate|log.Ltime),
		gameModel: &postgres.GameModel{DB: db},
	}

	log.Println("listening on port 4000")
	if err := http.ListenAndServe(":4000", s.loggingMiddleware(s.routes())); err != nil {
		return fmt.Errorf("error while serving: %v", err)
	}

	return nil
}

func openDB(host string, port int, user string, dbName string) (*sql.DB, error) {
	db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable", host, port, user, dbName))
	if err != nil {
		return nil, fmt.Errorf("error while opening connection to db: %w", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error while pinging db: %w", err)
	}

	return db, nil
}
