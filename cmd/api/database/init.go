package database

import "github.com/pressly/goose"

const (
	upCommand    = "up"
	sqlDirectory = "./sql"
)

// Init runs the migrations, needed to have a working Gira DB.
func Init(opts *DBOptions) error {
	db, err := NewDB(opts)
	if err != nil {
		return err
	}

	if err := goose.Run(upCommand, db, sqlDirectory); err != nil {
		return err
	}

	return nil
}
