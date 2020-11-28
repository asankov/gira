package postgres

import (
	"database/sql"
	"fmt"

	"github.com/asankov/gira/pkg/models"
	"github.com/lib/pq"
)

type FranchiseModel struct {
	DB *sql.DB
}

func (m *FranchiseModel) Insert(franchise *models.Franchise) (*models.Franchise, error) {
	row := m.DB.QueryRow(`INSERT INTO FRANCHISES (name, user_id) VALUES ($1, $2) RETURNING name, user_id`, franchise.Name, franchise.UserID)

	var f models.Franchise
	if err := row.Scan(&f.ID, &f.Name); err != nil {
		return nil, handleInsertFranchiseError(err)
	}

	return &f, nil
}

func handleInsertFranchiseError(err error) error {
	if err, ok := err.(*pq.Error); ok {
		if err.Constraint == "franchises_name_key" {
			return ErrNameAlreadyExists
		}
	}
	return fmt.Errorf("error while inserting record into the database: %w", err)
}

func (m *FranchiseModel) All(userID string) ([]*models.Franchise, error) {
	rows, err := m.DB.Query(`SELECT id, name FROM FRANCHISES f WHERE f.user_id = $1`, userID)
	if err != nil {
		return nil, fmt.Errorf("error while fetching franchises from the database: %w", err)
	}
	defer rows.Close()

	franchises := []*models.Franchise{}
	for rows.Next() {
		var franchise models.Franchise

		if err = rows.Scan(&franchise.ID, &franchise.Name); err != nil {
			return nil, fmt.Errorf("error while reading franchises from the database: %w", err)
		}

		franchises = append(franchises, &franchise)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error while reading franchises from the database: %w", err)
	}

	return franchises, nil
}
