package postgres

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/gira-games/api/pkg/models"
)

type FranchiseModel struct {
	DB *sql.DB
}

func (m *FranchiseModel) Insert(franchise *models.Franchise) (*models.Franchise, error) {
	if _, err := m.DB.Exec(`INSERT INTO FRANCHISES (name) VALUES ($1)`, franchise.Name); err != nil {
		if strings.Contains(err.Error(), `duplicate key value violates unique constraint "franchises_name_key"`) {
			return nil, ErrNameAlreadyExists
		}
		return nil, fmt.Errorf("error while inserting record into the database: %w", err)
	}

	var f models.Franchise
	if err := m.DB.QueryRow(`SELECT F.ID, F.NAME FROM FRANCHISES F WHERE F.NAME = $1`, franchise.Name).Scan(&f.ID, &f.Name); err != nil {
		return nil, fmt.Errorf("error while inserting record into the database: %w", err)
	}

	return &f, nil
}

func (m *FranchiseModel) All() ([]*models.Franchise, error) {
	rows, err := m.DB.Query(`SELECT id, name FROM FRANCHISES`)
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
