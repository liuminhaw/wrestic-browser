package restic

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"os/exec"
	"time"
)

const (
	awsAccessKeyIdEnv     string = "AWS_ACCESS_KEY_ID"
	awsSecretAccessKeyEnv string = "AWS_SECRET_ACCESS_KEY"
)

type S3Repository struct {
	Id              int
	Name            string
	Password        string
	Destination     string
	AccessKeyId     string
	SecretAccessKey string
	Region          string
	ConfigId        int
}

// Connect test if connection can be established using s3 bucket as backup repository
func (r *S3Repository) connect() error {
	os.Setenv(passwordEnv, r.Password)
	r.initCredential()

	commandArg := []string{"cat", "config", "-r", fmt.Sprintf("s3:s3.amazonaws.com/%s", r.Destination)}

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

func (r *S3Repository) newRepo(DB *sql.DB) error {
	// TODO: turn two query into transaction to ensure data is inserted atomically
	// TODO: encrypt password before inserting to database
	row := DB.QueryRow(`
		INSERT INTO "repositories" ("name", "destination", "password_enc", "type_id")
		VALUES ($1, $2, $3, (
			SELECT "id"
			FROM "repository_types"
			WHERE "name" = $4
		))		
		RETURNING ID;
	`, r.Name, r.Destination, r.Password, s3Type)
	err := row.Scan(&r.Id)
	if err != nil {
		return fmt.Errorf("create s3 repository: %w", err)
	}

	// TODO: encrypt access key and secret key before inserting to database
	row = DB.QueryRow(`
		INSERT INTO "s3_repository_configs" ("access_key_id_enc", "secret_access_key_enc", "region", "repository_id")		
		VALUES ($1, $2, $3, $4)
		RETURNING ID;
	`, r.AccessKeyId, r.SecretAccessKey, r.Region, r.Id)
	err = row.Scan(&r.ConfigId)
	if err != nil {
		return fmt.Errorf("create s3 repository config: %w", err)
	}

	return nil
}

func (r *S3Repository) initCredential() {
	os.Setenv(awsAccessKeyIdEnv, r.AccessKeyId)
	os.Setenv(awsSecretAccessKeyEnv, r.SecretAccessKey)
}
