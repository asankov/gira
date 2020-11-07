package database

import "github.com/pressly/goose"

const (
	upCommand = "up"
)

// Init runs the migrations, needed to have a working Gira DB.
// It gets the SQL files from the default `./sql` directory.
func Init(opts *DBOptions) error {
	return InitFromDirectory(opts, "./sql")
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
