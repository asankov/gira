package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	// to register PostreSQL driver
	_ "github.com/lib/pq"
)

type server struct {
	log     log.Logger
	handler http.Handler
}

func main() {
	if err := run(); err != nil {
		log.Panic("error while running server: " + err.Error())
	}
}

func run() error {
	s := &server{}

	db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable", "localhost", 5432, "antonsankov", "gira"))
	if err != nil {
		return fmt.Errorf("error while opening connection to db: %w", err)
	}
	if err := db.Ping(); err != nil {
		return fmt.Errorf("error while pinging db: %w", err)
	}
	defer db.Close()

	log.Println("listening on port 4000")
	if err := http.ListenAndServe(":4000", s.routes()); err != nil {
		return fmt.Errorf("error while serving: %v", err)
	}

	return nil
}
