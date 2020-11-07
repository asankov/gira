package database

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/sirupsen/logrus"
)

// DBOptions is the struct that encapsulates the fields needed to construct a database.
type DBOptions struct {
	Host   string
	Port   int
	User   string
	DBName string
	DBPass string
	UseSSL bool
}

var (
	// ErrDBHostMandatory is the error that indicates that the DB host is mandatory, but has not been set.
	ErrDBHostMandatory = errors.New(`db host is mandatory`)
	// ErrDBPortMandatory is the error that indicates that the DB port is mandatory, but has not been set.
	ErrDBPortMandatory = errors.New(`db port is mandatory`)
	// ErrDBUserMandatory is the error that indicates that the DB user is mandatory, but has not been set.
	ErrDBUserMandatory = errors.New(`db user is mandatory`)
	// ErrDBNameMandatory is the error that indicates that the DB name is mandatory, but has not been set.
	ErrDBNameMandatory = errors.New(`db name is mandatory`)
	// ErrDBPasswordMandatory is the error that indicates that the DB password is mandatory, but has not been set.
	ErrDBPasswordMandatory = errors.New(`db password is mandatory`)
)

// Validate validates the values of o, and error an error if any of the values is invalid
// If an error is returned it will be of type *multierror.Error and will contains information
// about all invalid fields.
func (o *DBOptions) Validate() error {
	var err *multierror.Error
	if o.Host == "" {
		err = multierror.Append(err, ErrDBHostMandatory)
	}
	if o.Port == 0 {
		err = multierror.Append(err, ErrDBPortMandatory)
	}
	if o.User == "" {
		err = multierror.Append(err, ErrDBUserMandatory)
	}
	if o.DBName == "" {
		err = multierror.Append(err, ErrDBNameMandatory)
	}
	if o.DBPass == "" {
		err = multierror.Append(err, ErrDBPasswordMandatory)
	}
	return err.ErrorOrNil()
}

// ConnectionString returns a connection string build from the values of o.
// This method will call Validate beforehand, and will propagate any error returned from Validate.
// The connections string will be in the format of PostgresSQL.
func (o *DBOptions) ConnectionString() (string, error) {
	if err := o.Validate(); err != nil {
		return "", err
	}

	connString := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s", o.Host, o.Port, o.User, o.DBName, o.DBPass)
	if !o.UseSSL {
		connString += " sslmode=disable"
	}

	return connString, nil
}

// NewDB builds a new *sql.DB from the passed options.
func NewDB(opts *DBOptions) (db *sql.DB, err error) {
	connString, err := opts.ConnectionString()
	if err != nil {
		return nil, err
	}

	logrus.Debugf("DB connection string: %s", connString)
	pings := 0
	for db, err = openDB(connString); err != nil; db, err = openDB(connString) {
		pings++
		time.Sleep(time.Duration(pings) * time.Second)
		logrus.Infof("retrying DB connection...%d\n", pings)
		if pings > 10 {
			return nil, err
		}
	}

	return db, nil
}

func openDB(connString string) (*sql.DB, error) {
	db, err := sql.Open("postgres", connString)
	if err != nil {
		return nil, fmt.Errorf("error while opening connection to db: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error while pinging db: %w", err)
	}

	return db, nil
}
