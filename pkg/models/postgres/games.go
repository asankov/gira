package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/asankov/gira/pkg/models"
)

var (
	// ErrNameAlreadyExists is the error that is returned, when a game with that name already exists in the database
	ErrNameAlreadyExists = errors.New("game with that name already exists in the database")
	// ErrNoRecord is returned when a game with that criteria does not exist in the database
	ErrNoRecord = errors.New("such game does not exist in the database")
)

// GameModel wraps an sql.DB connection pool.
type GameModel struct {
	DB *sql.DB
}

// Insert inserts the passed Game into the database.
// It returns the ID of the created game, or error if such occurred.
// If a game with the same name already exists, an ErrNameAlreadyExists is returned
func (m *GameModel) Insert(game *models.Game) (*models.Game, error) {
	if _, err := m.DB.Exec(`INSERT INTO GAMES (name) VALUES ($1)`, game.Name); err != nil {
		if strings.Contains(err.Error(), `duplicate key value violates unique constraint "games_name_key"`) {
			return nil, ErrNameAlreadyExists
		}
		return nil, fmt.Errorf("error while inserting record into the database: %w", err)
	}

	var g models.Game
	if err := m.DB.QueryRow(`SELECT * FROM GAMES G WHERE G.NAME = $1`, game.Name).Scan(&g.ID, &g.Name); err != nil {
		return nil, fmt.Errorf("error while inserting record into the database: %w", err)
	}

	return &g, nil
}

// Get fetches a Game by ID and returns that or an error if such occurred.
// If game with that ID is not present in the database, an ErrNoRecord is returned.
func (m *GameModel) Get(id string) (*models.Game, error) {
	var g models.Game

	if err := m.DB.QueryRow(`SELECT id, name FROM GAMES g WHERE g.id = $1`, id).Scan(&g.ID, &g.Name); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		}
		return nil, fmt.Errorf("error while fetching game from the database: %w", err)
	}

	return &g, nil
}

// All fetches all games from the database and returns them, or an error if such occurred.
func (m *GameModel) All() ([]*models.Game, error) {
	rows, err := m.DB.Query(`SELECT id, name FROM GAMES`)
	if err != nil {
		return nil, fmt.Errorf("error while fetching games from the database: %w", err)
	}
	defer rows.Close()

	var games []*models.Game
	for rows.Next() {
		var game models.Game

		if err = rows.Scan(&game.ID, &game.Name); err != nil {
			return nil, fmt.Errorf("error while reading games from the database: %w", err)
		}

		games = append(games, &game)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error while reading games from the database: %w", err)
	}

	return games, nil
}
