package restic

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"os/exec"
	"time"
)

type LocalRepository struct {
	Id          int
	Name        string
	Password    string
	Destination string
}

// Connect test if connection can be established using local directory as backup repository
func (r *LocalRepository) connect() error {
	os.Setenv(passwordEnv, r.Password)

	commandArg := []string{"cat", "config", "-r", r.Destination}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, resticCmd, commandArg...)
	_, err := cmd.Output()
	if ctx.Err() == context.DeadlineExceeded {
		return ErrConnectionTimeout
	}
	if err != nil {
		return fmt.Errorf("restic connect: %w", err)
	}

	return nil
}

func (r *LocalRepository) newRepo(DB *sql.DB) error {
	// TODO: encrypt password before inserting to database
	row := DB.QueryRow(`
		INSERT INTO "repositories" ("name", "destination", "password_enc", "type_id")
		VALUES ($1, $2, $3, (
			SELECT "id"
			FROM "repository_types"
			WHERE "name" = $4
		))		
		RETURNING ID;
	`, r.Name, r.Destination, r.Password, localType)
	err := row.Scan(&r.Id)
	if err != nil {
		return fmt.Errorf("create local repository: %w", err)
	}

	return nil
}
