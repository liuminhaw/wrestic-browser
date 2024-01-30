package restic

import (
	"database/sql"
	"fmt"
)

const (
	resticCmd   string = "restic"
	passwordEnv string = "RESTIC_PASSWORD"

	localType string = "local"
	s3Type    string = "s3"
	sftpType  string = "sftp"
)

type Repository interface {
	Connect() error
	// NewRepo() error
}

type RepositoryService struct {
	DB *sql.DB

	Repository Repository
}

// Types query and return all available repository types from respository_types table
func (service *RepositoryService) Types() ([]string, error) {
	rows, err := service.DB.Query(`
		SELECT name
		FROM repository_types
	`)
	if err != nil {
		return nil, fmt.Errorf("query repository types: %w", err)
	}
	defer rows.Close()

	var repoTypes []string
	for rows.Next() {
		var repoType string
		if err := rows.Scan(&repoType); err != nil {
			return nil, fmt.Errorf("query repository types: %w", err)
		}
		repoTypes = append(repoTypes, repoType)
	}
	if rows.Err() != nil {
		return nil, fmt.Errorf("query repository types: %w", err)
	}

	return repoTypes, nil
}
