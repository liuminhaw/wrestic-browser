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

// PemEnc is the encrypted PEM key used for SFTP authentication.
// PasswordEnc is the encrypted password used for SFTP authentication.
type SftpRepositoryEnc struct {
	PemEnc      string
	PasswordEnc string
}

// Id represents the unique identifier of the SFTP repository
// Name represents the name of the SFTP repository
// Password represents the password for the SFTP repository
// Destination represents the destination path for the SFTP repository
// User represents the username for the SFTP repository
// Host represents the host address for the SFTP repository
// Pem represents the PEM file path for the SFTP repository
// ConfigId represents the configuration ID for the SFTP repository
// Encryption represents the encryption settings for the SFTP repository
type SftpRepository struct {
	Id          int
	Name        string
	Password    string
	Destination string
	User        string
	Host        string
	Pem         string
	ConfigId    int
	Encryption  *SftpRepositoryEnc
}

// connect establishes a connection to the SFTP repository.
// It sets the password environment variable, creates a temporary PEM file,
// writes the PEM content to the file, sets the file permissions, and executes
// an SSH command to test the SFTP connection.
// If the connection test is successful, it executes a restic command to retrieve
// the repository configuration.
// Returns an error if any of the steps fail.
func (r *SftpRepository) connect() error {
	os.Setenv(passwordEnv, r.Password)

	// Create pem temporary pem file
	f, err := os.CreateTemp("", "wrestic-brw-pem")
	if err != nil {
		return fmt.Errorf("restic connect: create temp: %w", err)
	}
	defer os.Remove(f.Name())

	tempFilename := f.Name()
	fmt.Printf("Temp file name: %s\n", tempFilename)

	if _, err := f.WriteString(r.Pem); err != nil {
		return fmt.Errorf("restic connect: write pem: %w", err)
	}
	if err := f.Chmod(0600); err != nil {
		return fmt.Errorf("restic connect: file chmod: %w", err)
	}
	if err := f.Close(); err != nil {
		return fmt.Errorf("restic connect: close file: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// sftp connection test
	commandArg := []string{
		"-i",
		tempFilename,
		"-o",
		"PasswordAuthentication=no",
		"-o",
		"StrictHostKeyChecking=no",
		"-o",
		"ServerAliveInterval=60",
		"-o",
		"ServerAliveCountMax=240",
		"-o",
		"BatchMode=yes",
		fmt.Sprintf("%s@%s", r.User, r.Host),
		"ls",
		r.Destination,
		">",
		"/dev/null",
	}
	cmd := exec.CommandContext(ctx, "ssh", commandArg...)
	output, err := cmd.CombinedOutput()
	if ctx.Err() == context.DeadlineExceeded {
		return ErrConnectionTimeout
	}
	if err != nil {
		return fmt.Errorf("restic connect: ssh test: %s: %w", output, err)
	}

	commandArg = []string{
		"cat",
		"config",
		"-r",
		fmt.Sprintf("sftp::%s", r.Destination),
		"-o",
		fmt.Sprintf("sftp.command=ssh %s@%s -o PasswordAuthentication=no -o StrictHostKeyChecking=no -o ServerAliveInterval=60 -o ServerAliveCountMax=240 -o BatchMode=yes -i %s -T -s sftp", r.User, r.Host, tempFilename),
	}

	cmd = exec.CommandContext(ctx, resticCmd, commandArg...)

	output, err = cmd.CombinedOutput()
	if ctx.Err() == context.DeadlineExceeded {
		return ErrConnectionTimeout
	}
	if err != nil {
		return fmt.Errorf("restic connect: %s: %w", output, err)
	}

	return nil
}

// newRepo creates a new SFTP repository in the database.
// It takes a DB connection as input and returns an error if any.
// The function inserts the repository details into the "repositories" table,
// and the SFTP repository configuration into the "sftp_repository_configs" table.
// It uses a transaction to ensure atomicity and rolls back the transaction if any error occurs.
func (r *SftpRepository) newRepo(DB *sql.DB, userId int) error {
	tx, err := DB.Begin()
	if err != nil {
		return fmt.Errorf("create sftp repository: new transaction: %w", err)
	}

	row := tx.QueryRow(`
		INSERT INTO "repositories" ("name", "destination", "password_enc", "type_id", "owner_id")
		VALUES ($1, $2, $3, (
			SELECT "id"
			FROM "repository_types"
			WHERE "name" = $4
		), $5)
		RETURNING ID;
	`, r.Name, r.Destination, r.Encryption.PasswordEnc, sftpType, userId)
	err = row.Scan(&r.Id)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("create sftp repository: transaction rollback: %w", rbErr)
		}
		return fmt.Errorf("create sftp repository: %w", err)
	}

	row = tx.QueryRow(`
		INSERT INTO "sftp_repository_configs" ("user", "host", "pem_enc", "repository_id")		
		VALUES ($1, $2, $3, $4)
		RETURNING ID;
	`, r.User, r.Host, r.Encryption.PemEnc, r.Id)
	err = row.Scan(&r.ConfigId)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("create sftp repository: transaction rollback: %w", rbErr)
		}
		return fmt.Errorf("create sftp repository config: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("create sftp repository: transaction commit: %w", err)
	}

	return nil
}

// GenEnc generates encryption for the SftpRepository using the provided encryption key.
// It encrypts the repository's PEM and password using the given encryption key.
// The encrypted values are stored in the repository's Encryption struct.
// If an error occurs during encryption, it returns an error with a descriptive message.
func (r *SftpRepository) GenEnc(encKey [32]byte) error {
	enc, err := encryptor.Encrypt([]byte(r.Pem), encKey)
	if err != nil {
		return fmt.Errorf("gen sftp repository encryption: pem: %w", err)
	}
	r.Encryption.PemEnc = enc

	enc, err = encryptor.Encrypt([]byte(r.Password), encKey)
	if err != nil {
		return fmt.Errorf("gen sftp repository encryption: password: %w", err)
	}
	r.Encryption.PasswordEnc = enc

	return nil
}
