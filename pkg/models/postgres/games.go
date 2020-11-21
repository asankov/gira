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
	ErrNameAlreadyExists = errors.New("model with that name already exists in the database")
	// ErrNoRecord is returned when a game with that criteria does not exist in the database
	ErrNoRecord = errors.New("such model does not exist in the database")
)

// GameModel wraps an sql.DB connection pool.
type GameModel struct {
	DB *sql.DB
}

// Insert inserts the passed Game into the database.
// It returns the ID of the created game, or error if such occurred.
// If a game with the same name already exists, an ErrNameAlreadyExists is returned
func (m *GameModel) Insert(game *models.Game) (*models.Game, error) {
	var err error
	if game.FranchiseID == "" {
		_, err = m.DB.Exec(`INSERT INTO GAMES (name, user_id) VALUES ($1, $2)`, game.Name, game.UserID)
	} else {
		_, err = m.DB.Exec(`INSERT INTO GAMES (name, user_id, franchise_id) VALUES ($1, $2, $3)`, game.Name, game.UserID, game.FranchiseID)
	}

	// TODO: proper error handling
	if err != nil {
		if strings.Contains(err.Error(), `duplicate key value violates unique constraint "games_name_key"`) {
			return nil, ErrNameAlreadyExists
		}
		return nil, fmt.Errorf("error while inserting record into the database: %w", err)
	}

	// TODO: use returning
	var g models.Game
	var fID sql.NullString
	if err := m.DB.QueryRow(`SELECT G.ID, G.NAME, G.FRANCHISE_ID FROM GAMES G WHERE G.NAME = $1 AND G.USER_ID = $2`, game.Name, game.UserID).Scan(&g.ID, &g.Name, &fID); err != nil {
		return nil, fmt.Errorf("error while inserting record into the database: %w", err)
	}
	g.FranchiseID = fID.String

	return &g, nil
}

// Get fetches a Game by ID and returns that or an error if such occurred.
// If game with that ID is not present in the database, an ErrNoRecord is returned.
func (m *GameModel) Get(id string) (*models.Game, error) {
	var g models.Game

	if err := m.DB.QueryRow(`SELECT g.id, g.name, g.franchise_id, f.name FROM GAMES g WHERE g.id = $1 JOIN FRANCHISES f ON f.id = g.franchise_id`, id).Scan(&g.ID, &g.Name, &g.FranchiseID, &g.Franchise); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		}
		return nil, fmt.Errorf("error while fetching game from the database: %w", err)
	}

	return &g, nil
}

// AllForUser fetches all games for the given user from the database and returns them, or an error if such occurred.
func (m *GameModel) AllForUser(userID string) ([]*models.Game, error) {
	rows, err := m.DB.Query(`
	SELECT 
		g.id, 
		g.name, 
		g.franchise_id, 
		f.name AS frachise_name,
		g.status,
		g.current_progress,
		g.final_progress
	FROM GAMES g 
		LEFT JOIN FRANCHISES f ON f.id = g.franchise_id 
	WHERE g.user_id = $1`, userID)
	if err != nil {
		return nil, fmt.Errorf("error while fetching games from the database: %w", err)
	}
	defer rows.Close()

	games := []*models.Game{}
	for rows.Next() {
		game := models.Game{Progress: &models.GameProgress{}}

		var fID, fName sql.NullString
		if err = rows.Scan(&game.ID, &game.Name, &fID, &fName, &game.Status, &game.Progress.Current, &game.Progress.Final); err != nil {
			return nil, fmt.Errorf("error while reading games from the database: %w", err)
		}
		game.Franchise = fName.String
		game.FranchiseID = fID.String

		games = append(games, &game)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error while reading games from the database: %w", err)
	}

	return games, nil
}

func (m *GameModel) DeleteGame(userID, gameID string) error {
	if _, err := m.DB.Exec(`DELETE FROM GAMES U WHERE U.id = $1 AND U.user_id = $2`, userID, gameID); err != nil {
		return err
	}
	return nil
}

func (m *GameModel) ChangeGameStatus(userID, gameID string, status models.Status) error {
	if _, err := m.DB.Exec("UPDATE GAMES SET status = $1 WHERE id = $2 AND user_id = $3", status, gameID, userID); err != nil {
		return fmt.Errorf("error while updating game status: %w", err)
	}
	return nil
}

func (m *GameModel) ChangeGameProgress(userID, gameID string, progress *models.GameProgress) error {
	if _, err := m.DB.Exec("UPDATE GAMES SET current_progress =  $1, final_progress = $2 WHERE id = $3 AND user_id = $4", progress.Current, progress.Final, gameID, userID); err != nil {
		return fmt.Errorf("error while updating game progress: %w", err)
	}
	return nil
}
