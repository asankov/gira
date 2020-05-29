package postgres

import (
	"database/sql"
	"fmt"

	"github.com/asankov/gira/pkg/models"
)

type UserGamesModel struct {
	DB *sql.DB
}

func (m *UserGamesModel) LinkGameToUser(userID, gameID string) error {
	if _, err := m.DB.Exec(`INSERT INTO USER_GAMES(user_id, game_id) VALUES ($1, $2)`, userID, gameID); err != nil {
		return err
	}
	return nil
}

func (m *UserGamesModel) GetUserGames(userID string) ([]*models.Game, error) {
	rows, err := m.DB.Query(`SELECT id, name FROM GAMES WHERE id IN (SELECT game_id FROM USER_GAMES WHERE user_id = $1)`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	games := []*models.Game{}
	for rows.Next() {
		var game models.Game

		if err = rows.Scan(&game.ID, &game.Name); err != nil {
			return nil, fmt.Errorf("error while reading games from the database: %w", err)
		}

		games = append(games, &game)
	}

	return games, nil
}
