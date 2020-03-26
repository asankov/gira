package postgres

import (
	"database/sql"
	"fmt"

	"github.com/asankov/gira/pkg/models"
)

// GameModel wraps an sql.DB connection pool.
type GameModel struct {
	DB *sql.DB
}

// Insert inserts the passed Game into the database.
// It returns the ID of the created game, or error if such occurred.
func (m *GameModel) Insert(game *models.Game) (string, error) {
	res, err := m.DB.Exec(`INSERT INTO GAMES (name) VALUES ($1)`, game.Name)
	if err != nil {
		return "", fmt.Errorf("error while inserting record into the database: %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return "", fmt.Errorf("error while inserting record into the database: %w", err)
	}
	return string(id), nil
}

// Get fetches a Game by ID and returns that or an error if such occurred.
func (m *GameModel) Get(id string) (*models.Game, error) {
	return nil, nil
}

// All fetches all games from the database and returns them, or an error if such occurred.
func (m *GameModel) All() ([]*models.Game, error) {
	return nil, nil
}
