package restic

import (
	"database/sql"
	"fmt"

	"github.com/jackc/pgtype"
)

type RepositoryStatusService struct {
	DB *sql.DB
}

type RepositoryStatus struct {
	Name       string
	Status     string
	LastBackup string
	Owner      string
}

// list returns the status of each repository.
// It returns a slice of RepositoryStatus or an error if the query fails.
func (service *RepositoryStatusService) List() ([]RepositoryStatus, error) {
	// Query the database to get the status of each repository
	// This get all the repositories which is for admin user,
	// TODO: normal user will only get the repositories which is owned by the user
	rows, err := service.DB.Query(`
        SELECT
            "repositories"."name",
            "repository_settings"."backup_status",
            "repository_settings"."recent_snapshot",
            "users"."username"
        FROM "repositories"
        INNER JOIN "repository_settings" 
            ON "repositories"."id" = "repository_settings"."repository_id"
        INNER JOIN "users" ON "repositories"."owner_id" = "users"."id"
    `)
	if err != nil {
		return nil, fmt.Errorf("repository status list: query repository status: %w", err)
	}
	defer rows.Close()

	var statuses []RepositoryStatus
	for rows.Next() {
		var status RepositoryStatus
		var lastBackupTime pgtype.Timestamptz
		err := rows.Scan(&status.Name, &status.Status, &lastBackupTime, &status.Owner)
		if err != nil {
			return nil, fmt.Errorf("repository status list: read status: %w", err)
		}

		// Format the last backup time
		if lastBackupTime.Status == pgtype.Present {
			status.LastBackup = lastBackupTime.Time.Format("2006-01-02 15:04:05")
		} else {
			status.LastBackup = "N/A"
		}

		statuses = append(statuses, status)
	}
	if rows.Err() != nil {
		return nil, fmt.Errorf("repository status list: read status: %w", err)
	}

	return statuses, nil
}
