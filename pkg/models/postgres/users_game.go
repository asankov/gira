package postgres

import (
	"database/sql"
	"fmt"

	"github.com/gira-games/api/pkg/models"
)

type UserGamesModel struct {
	DB *sql.DB
}

func (m *UserGamesModel) LinkGameToUser(userID, gameID string, progress *models.UserGameProgress) error {
	if _, err := m.DB.Exec(`INSERT INTO USER_GAMES(user_id, game_id, current_progress, final_progress) VALUES ($1, $2, $3, $4)`, userID, gameID, progress.Current, progress.Final); err != nil {
		return err
	}
	return nil
}

func (m *UserGamesModel) DeleteUserGame(userGameID string) error {
	if _, err := m.DB.Exec(`DELETE FROM USER_GAMES U WHERE U.id = $1`, userGameID); err != nil {
		return err
	}
	return nil
}

func (m *UserGamesModel) GetUserGames(userID string) ([]*models.UserGame, error) {
	rows, err := m.DB.Query(`
		SELECT
			ug.id,
			g.id AS user_game_id,
			g.name,
			ug.status,
			ug.current_progress,
			ug.final_progress,
			f.name
		FROM USER_GAMES ug
			JOIN GAMES g ON ug.game_id = g.id
			LEFT JOIN FRANCHISES f ON f.id = g.franchise_id
		WHERE ug.user_id = $1`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	userGames := []*models.UserGame{}
	for rows.Next() {
		var userGame = models.UserGame{
			Game:     &models.Game{},
			Progress: &models.UserGameProgress{},
		}

		var fName sql.NullString
		if err = rows.Scan(&userGame.ID, &userGame.Game.ID, &userGame.Game.Name, &userGame.Status, &userGame.Progress.Current, &userGame.Progress.Final, &fName); err != nil {
			return nil, fmt.Errorf("error while reading games from the database: %w", err)
		}
		userGame.Game.Franchise = fName.String

		userGames = append(userGames, &userGame)
	}

	return userGames, nil
}

func (m *UserGamesModel) GetUserGamesGrouped(userID string) (map[models.Status][]*models.UserGame, error) {

	games, err := m.GetUserGames(userID)
	if err != nil {
		return nil, err
	}

	// TODO: this is really stupid. all of this should be implemented via SQL group by
	// and ideally the statuses should be customizable
	todo, inProgress, done := []*models.UserGame{}, []*models.UserGame{}, []*models.UserGame{}
	for _, game := range games {
		if game.Status == models.StatusTODO {
			todo = append(todo, game)
		} else if game.Status == models.StatusInProgress {
			inProgress = append(inProgress, game)
		} else if game.Status == models.StatusDone {
			done = append(done, game)
		}
	}

	return map[models.Status][]*models.UserGame{
		models.StatusTODO:       todo,
		models.StatusInProgress: inProgress,
		models.StatusDone:       done,
	}, nil
}

func (m *UserGamesModel) ChangeGameStatus(userID, userGameID string, status models.Status) error {
	if _, err := m.DB.Exec("UPDATE USER_GAMES SET status =  $1 WHERE id = $2 AND user_id = $3", status, userGameID, userID); err != nil {
		return fmt.Errorf("error while updating game status: %w", err)
	}
	return nil
}

func (m *UserGamesModel) ChangeGameProgress(userID, userGameID string, progress *models.UserGameProgress) error {
	if _, err := m.DB.Exec("UPDATE USER_GAMES SET current_progress =  $1, final_progress = $2 WHERE id = $3 AND user_id = $4", progress.Current, progress.Final, userGameID, userID); err != nil {
		return fmt.Errorf("error while updating game progress: %w", err)
	}
	return nil
}

func (m *UserGamesModel) GetAvailableGamesFor(userID string) ([]*models.Game, error) {
	rows, err := m.DB.Query(`
		SELECT 
			g.id, 
			g.name, 
			f.name AS franchise_name 
		FROM games g 
			LEFT JOIN FRANCHISES f ON f.id = g.franchise_id 
			WHERE g.id NOT IN (SELECT game_id FROM user_games u WHERE u.user_id = $1) AND g.user_id = $1`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	games := []*models.Game{}
	for rows.Next() {
		var game = models.Game{}

		var fName sql.NullString
		if err = rows.Scan(&game.ID, &game.Name, &fName); err != nil {
			return nil, fmt.Errorf("error while reading games from the database: %w", err)
		}
		game.Franchise = fName.String

		games = append(games, &game)
	}

	return games, nil

}
