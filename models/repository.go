package models

import (
	"database/sql"
	"fmt"
)

type RepositoryTypes struct {
	Name string
}

type RepositoryService struct {
	DB *sql.DB
}

func (service *RepositoryService) Types() ([]RepositoryTypes, error) {
	rows, err := service.DB.Query(`
		SELECT name
		FROM repository_types
	`)
	if err != nil {
		return nil, fmt.Errorf("query repository types: %w", err)
	}
	defer rows.Close()

	var repoTypes []RepositoryTypes
	for rows.Next() {
		repoType := RepositoryTypes{}
		if err := rows.Scan(&repoType.Name); err != nil {
			return nil, fmt.Errorf("query repository types: %w", err)
		}
		repoTypes = append(repoTypes, repoType)
	}
	if rows.Err() != nil {
		return nil, fmt.Errorf("query repository types: %w", err)
	}

	return repoTypes, nil
}
