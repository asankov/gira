package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/asankov/gira/pkg/models"
	"golang.org/x/crypto/bcrypt"
)

var (
	// ErrEmailAlreadyExists is returned when a user with the same email already exists in the database
	ErrEmailAlreadyExists = errors.New("user with the same email already exists")
	// ErrUsernameAlreadyExists is returned when a user with the same username already exists in the database
	ErrUsernameAlreadyExists = errors.New("user with the same username already")
)

// UserModel wraps a DB connection pool.
type UserModel struct {
	DB *sql.DB
}

// Insert inserts a new user with the given parameters into the database
// and return the created user ot an error if such occurred.
func (m *UserModel) Insert(user *models.User) (*models.User, error) {
	// TODO: hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
	if err != nil {
		return nil, fmt.Errorf("error while hashing password: %w", err)
	}

	if _, err := m.DB.Exec("INSERT INTO USERS (username, email, hashed_password) VALUES ($1, $2, $3)", user.Username, user.Email, hash); err != nil {
		return nil, handleInsertError(err)
	}

	usr := models.User{}
	if err := m.DB.QueryRow("SELECT id, username, email FROM USERS U WHERE U.USERNAME = $1 AND U.EMAIL = $2", user.Username, user.Email).Scan(&usr.ID, &usr.Username, &usr.Email); err != nil {
		return nil, fmt.Errorf("error while fetching user from the database: %w", err)
	}

	return &usr, nil
}

func handleInsertError(err error) error {
	if strings.Contains(err.Error(), `duplicate key value violates unique constraint "users_uc_email"`) {
		return ErrEmailAlreadyExists
	}
	if strings.Contains(err.Error(), `duplicate key value violates unique constraint "users_uc_username"`) {
		return ErrUsernameAlreadyExists
	}
	return fmt.Errorf("error while inserting user into the database: %w", err)
}

// Authenticate authenticates a use with these credentials
// and returns the user or an error if such occurred.
func (m *UserModel) Authenticate(email, password string) (*models.User, error) {
	return nil, nil
}

// Get fetches the user with the given ID from the database
// and returns the user or an error if such occurred.
func (m *UserModel) Get(id string) (*models.User, error) {
	return nil, nil
}
