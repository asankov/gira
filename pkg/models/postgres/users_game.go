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
	rows, err := m.DB.Query(`SELECT g.id, g.name, ug.status FROM USER_GAMES ug JOIN GAMES g ON ug.game_id = g.id where user_id = $1`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	games := []*models.Game{}
	for rows.Next() {
		var game models.Game

		if err = rows.Scan(&game.ID, &game.Name, &game.Status); err != nil {
			return nil, fmt.Errorf("error while reading games from the database: %w", err)
		}

		games = append(games, &game)
	}

	return games, nil
}

func (m *UserGamesModel) GetUserGamesGrouped(userID string) (map[models.Status][]*models.Game, error) {

	games, err := m.GetUserGames(userID)
	if err != nil {
		return nil, err
	}

	// TODO: this is really stupid. all of this should be implemented via SQL group by
	// and ideally the statuses should be customizable
	todo, inProgress, done := []*models.Game{}, []*models.Game{}, []*models.Game{}
	for _, game := range games {
		if game.Status == models.StatusTODO {
			todo = append(todo, game)
		} else if game.Status == models.StatusInProgress {
			inProgress = append(inProgress, game)
		} else if game.Status == models.StatusDone {
			done = append(done, game)
		}
	}

	return map[models.Status][]*models.Game{
		models.StatusTODO:       todo,
		models.StatusInProgress: inProgress,
		models.StatusDone:       done,
	}, nil
}

func (m *UserGamesModel) ChangeGameStatus(userID, gameID string, status models.Status) error {
	if _, err := m.DB.Exec("UPDATE USER_GAMES SET status =  $1 WHERE user_id = $2 AND game_id = $3", status, userID, gameID); err != nil {
		return fmt.Errorf("error while updating game status: %w", err)
	}
	return nil
}
