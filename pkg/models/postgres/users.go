package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/gira-games/api/pkg/models"
	"golang.org/x/crypto/bcrypt"
)

var (
	// ErrEmailAlreadyExists is returned when a user with the same email already exists in the database
	ErrEmailAlreadyExists = errors.New("user with the same email already exists")
	// ErrUsernameAlreadyExists is returned when a user with the same username already exists in the database
	ErrUsernameAlreadyExists = errors.New("user with the same username already exists")
	// ErrWrongPassword is returned when the given password does not match the user password
	ErrWrongPassword = errors.New("the given password does not match the user password")
)

// UserModel wraps a DB connection pool.
type UserModel struct {
	DB *sql.DB
}

// Insert inserts a new user with the given parameters into the database
// and return the created user ot an error if such occurred.
func (m *UserModel) Insert(user *models.User) (*models.User, error) {
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
	usr := models.User{}
	if err := m.DB.QueryRow("SELECT id, username, email, hashed_password FROM USERS U WHERE U.EMAIL = $1", email).Scan(&usr.ID, &usr.Username, &usr.Email, &usr.HashedPassword); err != nil {
		return nil, fmt.Errorf("error while fetching user from the database: %w", err)
	}
	if err := bcrypt.CompareHashAndPassword(usr.HashedPassword, []byte(password)); err != nil {
		return nil, ErrWrongPassword
	}
	return &usr, nil
}

// AssociateTokenWithUser associated the given token with the given userID
func (m *UserModel) AssociateTokenWithUser(userID, token string) error {
	if _, err := m.DB.Exec("INSERT INTO user_tokens (user_id, token) VALUES ($1, $2)", userID, token); err != nil {
		// TODO: better error handling
		return fmt.Errorf("error while inserting token into database: %w", err)
	}
	return nil
}

// InvalidateToken deleted the token from the database, making it invalid
func (m *UserModel) InvalidateToken(userID, token string) error {
	if _, err := m.DB.Exec("DELETE FROM user_tokens WHERE user_id = $1 AND token = $2", userID, token); err != nil {
		// TODO: better error handling
		return fmt.Errorf("error while deleting token from the database: %w", err)
	}
	return nil
}

// GetUserByToken returns the user, associated with the token passed to the method
func (m *UserModel) GetUserByToken(token string) (*models.User, error) {
	var usr models.User
	if err := m.DB.QueryRow("SELECT id, username, email FROM USERS U WHERE id = (SELECT user_id FROM user_tokens WHERE token = $1)", token).Scan(&usr.ID, &usr.Username, &usr.Email); err != nil {
		return nil, fmt.Errorf("error while looking up user: %w", err)
	}
	return &usr, nil
}
