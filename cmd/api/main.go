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
	s := &server{}

	db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable", "localhost", 5432, "antonsankov", "gira"))
	if err != nil {
		panic("err open: " + err.Error())
	}
	if err := db.Ping(); err != nil {
		panic("err ping: " + err.Error())
	}

	log.Println("listening on port 4000")
	if err := http.ListenAndServe(":4000", s.routes()); err != nil {
		log.Fatalf("error while serving: %v", err)
	}
}
