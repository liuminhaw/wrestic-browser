package restic

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/liuminhaw/wrestic-brw/utils/encryptor"
)

// LocalRepositoryEnc represents a local repository with encrypted password.
type LocalRepositoryEnc struct {
	PasswordEnc string
}

// Id: Unique identifier for the local repository
// Name: Name of the local repository
// Password: Password for the local repository
// PasswordEnc: Encrypted password for the local repository
// Destination: Destination path for the local repository
type LocalRepository struct {
	Id          int
	Name        string
	Password    string
	Destination string
	Encryption  *LocalRepositoryEnc
}

// connect establishes a connection to the local repository.
// It sets the password environment variable and executes the restic command to check the connection.
// If the connection times out, it returns ErrConnectionTimeout.
// If there is an error executing the restic command, it returns an error with the message "restic connect: <original error>".
// Otherwise, it returns nil.
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

// newRepo creates a new local repository in the database.
// It inserts the repository details into the "repositories" table,
// including the name, destination, encrypted password, and type ID.
// The type ID is obtained by querying the "repository_types" table
// based on the provided repository type name.
// The ID of the newly created repository is returned.
// If an error occurs during the database operation, it is wrapped
// and returned as an error.
func (r *LocalRepository) newRepo(DB *sql.DB, userId int) error {
	row := DB.QueryRow(`
		INSERT INTO "repositories" ("name", "destination", "password_enc", "type_id", "owner_id")
		VALUES ($1, $2, $3, (
			SELECT "id"
			FROM "repository_types"
			WHERE "name" = $4
		), $5)		
		RETURNING ID;
	`, r.Name, r.Destination, r.Encryption.PasswordEnc, localType, userId)
	err := row.Scan(&r.Id)
	if err != nil {
		return fmt.Errorf("create local repository: %w", err)
	}

	return nil
}

// GenEnc generates the encryption for the local repository using the provided key.
// It encrypts the repository password using the key and stores the encrypted value in the Encryption.PasswordEnc field.
// If an error occurs during encryption, it returns an error with a descriptive message.
func (r *LocalRepository) GenEnc(key [32]byte) error {
	enc, err := encryptor.Encrypt([]byte(r.Password), key)
	if err != nil {
		return fmt.Errorf("gen local repository encryption: %w", err)
	}

	r.Encryption.PasswordEnc = enc
	return nil
}
