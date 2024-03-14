package restic

import (
	"database/sql"
	"fmt"
	"time"
)

const (
	resticCmd   string = "restic"
	passwordEnv string = "RESTIC_PASSWORD"

	localType string = "local"
	s3Type    string = "s3"
	sftpType  string = "sftp"

	defaultSnapshotCheckInterval = 1 * 24 * time.Hour
)

type Repository interface {
	connect() error
	newRepo(*sql.DB, int) error
	GenEnc([32]byte) error
}

// RepositoryService represents a service that interacts with a repository.
type RepositoryService struct {
	DB         *sql.DB    // DB is the database connection.
	EncKey     [32]byte   // EncKey is the encryption key used for the repository.
	Repository Repository // Repository is the underlying repository.
}

// Types returns a list of repository types.
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

// Connect establishes a connection to the repository.
// It returns an error if the connection fails.
func (service *RepositoryService) Connect() error {
	if err := service.Repository.connect(); err != nil {
		return fmt.Errorf("repository connect failed: %w", err)
	}

	return nil
}

// Create creates a new repository.
// It calls the newRepo method of the RepositoryService to initialize the repository.
// If an error occurs during the creation process, it returns an error with a formatted message.
func (service *RepositoryService) Create(userId int) error {
	if err := service.Repository.newRepo(service.DB, userId); err != nil {
		return fmt.Errorf("create repository: %w", err)
	}

	return nil
}
