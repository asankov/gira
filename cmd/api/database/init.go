package database

import (
	"errors"
	"runtime"
	"strings"

	"github.com/pressly/goose"
)

const (
	upCommand = "up"
	// currentFilename points to the name of the current file.
	// its purpose is to be replaced by `sql/` to infer the path
	// to the SQL migrations.
	// it's important to change this if the `Directory` function is moved
	// into another file
	currentFilename = "cmd/api/database/init.go"
)

// Init runs the migrations, needed to have a working Gira DB.
// It gets the SQL files from the default `./sql` directory.
func Init(opts *DBOptions) error {
	dir, err := MigrationsDirectory()
	if err != nil {
		return err
	}
	return InitFromDirectory(opts, dir)
}

// InitFromDirectory runs the migrations, needed to have a working Gira DB.
// It gets the SQL files from the provided directory.
func InitFromDirectory(opts *DBOptions, sqlDirectory string) error {
	db, err := NewDB(opts)
	if err != nil {
		return err
	}

	if err := goose.Run(upCommand, db, sqlDirectory); err != nil {
		return err
	}

	return nil
}

// MigrationsDirectory returns the path to the directory containing the SQL migrations.
// It will return an absolute path that can be passed directly to `InitFromDirectory`.
func MigrationsDirectory() (string, error) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", errors.New("not able to recover the current directory")
	}
	return strings.ReplaceAll(filename, currentFilename, "sql"), nil
}
