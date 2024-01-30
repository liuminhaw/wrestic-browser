package restic

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"os/exec"
	"time"
)

type SftpRepository struct {
	Id          int
	Name        string
	Password    string
	Destination string
	User        string
	Host        string
	Pem         string
	ConfigId    int
}

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
		fmt.Sprintf("%s", tempFilename),
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
		fmt.Sprintf("%s", r.Destination),
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

func (r *SftpRepository) newRepo(DB *sql.DB) error {
	// TODO: turn two queries into transaction to ensure data is inserted atomically
	// TODO: encrypt password before inserting to database
	row := DB.QueryRow(`
		INSERT INTO "repositories" ("name", "destination", "password_enc", "type_id")
		VALUES ($1, $2, $3, (
			SELECT "id"
			FROM "repository_types"
			WHERE "name" = $4
		))
		RETURNING ID;
	`, r.Name, r.Destination, r.Password, sftpType)
	err := row.Scan(&r.Id)
	if err != nil {
		return fmt.Errorf("create sftp repository: %w", err)
	}

	// TODO: encrypt pem before inserting to database
	row = DB.QueryRow(`
		INSERT INTO "sftp_repository_configs" ("user", "host", "pem_enc", "repository_id")		
		VALUES ($1, $2, $3, $4)
		RETURNING ID;
	`, r.User, r.Host, r.Pem, r.Id)
	err = row.Scan(&r.ConfigId)
	if err != nil {
		return fmt.Errorf("create sftp repository config: %w", err)
	}

	return nil
}
