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

const (
	awsAccessKeyIdEnv     string = "AWS_ACCESS_KEY_ID"
	awsSecretAccessKeyEnv string = "AWS_SECRET_ACCESS_KEY"
)

// Encrypted access key ID for S3 repository
// Encrypted secret access key for S3 repository
// Encrypted password for S3 repository
type S3RepositoryEnc struct {
	AccessKeyIdEnc     string
	SecretAccessKeyEnc string
	PasswordEnc        string
}

// Id represents the unique identifier of the S3 repository.
// Name represents the name of the S3 repository.
// Password represents the password for accessing the S3 repository.
// Destination represents the target location where the S3 repository is stored.
// AccessKeyId represents the access key ID for authenticating with the S3 repository.
// SecretAccessKey represents the secret access key for authenticating with the S3 repository.
// Region represents the AWS region where the S3 repository is located.
// ConfigId represents the configuration identifier associated with the S3 repository.
// Encryption represents the encryption settings for the S3 repository.
type S3Repository struct {
	Id              int
	Name            string
	Password        string
	Destination     string
	AccessKeyId     string
	SecretAccessKey string
	Region          string
	ConfigId        int
	Encryption      *S3RepositoryEnc
}

// connect establishes a connection to the S3 repository.
// It sets the password environment variable, initializes the credentials,
// and executes a restic command to check the connection.
// If the connection times out, it returns ErrConnectionTimeout.
// If there is an error during the connection, it returns an error with the message "restic connect: <original error>".
// Otherwise, it returns nil.
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

// newRepo creates a new S3 repository in the database.
// It takes a DB connection as input and returns an error if any.
// The function inserts the repository details into the "repositories" table
// and the S3 repository configuration into the "s3_repository_configs" table.
// It uses a transaction to ensure atomicity and rolls back the transaction if any error occurs.
func (r *S3Repository) newRepo(DB *sql.DB, userId int) error {
	tx, err := DB.Begin()
	if err != nil {
		return fmt.Errorf("create s3 repository: new transaction: %w", err)
	}

	row := tx.QueryRow(`
		INSERT INTO "repositories" ("name", "destination", "password_enc", "type_id", "owner_id")
		VALUES ($1, $2, $3, (
			SELECT "id"
			FROM "repository_types"
			WHERE "name" = $4
		), $5)	
		RETURNING ID;
	`, r.Name, r.Destination, r.Encryption.PasswordEnc, s3Type, userId)
	err = row.Scan(&r.Id)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("create s3 repository: transaction rollback: %w", rbErr)
		}
		return fmt.Errorf("create s3 repository: %w", err)
	}

	row = tx.QueryRow(`
		INSERT INTO "s3_repository_configs" ("access_key_id_enc", "secret_access_key_enc", "region", "repository_id")	
		VALUES ($1, $2, $3, $4)
		RETURNING ID;
	`, r.Encryption.AccessKeyIdEnc, r.Encryption.SecretAccessKeyEnc, r.Region, r.Id)
	err = row.Scan(&r.ConfigId)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("create s3 repository: transaction rollback: %w", rbErr)
		}
		return fmt.Errorf("create s3 repository config: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("create s3 repository: transaction commit: %w", err)
	}

	return nil
}

// GenEnc generates encrypted versions of the S3 repository's access key ID, secret access key, and password.
// It takes an encryption key as input and encrypts the values using the encryptor package.
// The encrypted values are then stored in the Encryption struct of the S3Repository.
// If any encryption operation fails, an error is returned.
func (r *S3Repository) GenEnc(encKey [32]byte) error {
	enc, err := encryptor.Encrypt([]byte(r.AccessKeyId), encKey)
	if err != nil {
		return fmt.Errorf("gen s3 repository encryption: access key id: %w", err)
	}
	r.Encryption.AccessKeyIdEnc = enc

	enc, err = encryptor.Encrypt([]byte(r.SecretAccessKey), encKey)
	if err != nil {
		return fmt.Errorf("gen s3 repository encryption: secret access key: %w", err)
	}
	r.Encryption.SecretAccessKeyEnc = enc

	enc, err = encryptor.Encrypt([]byte(r.Password), encKey)
	if err != nil {
		return fmt.Errorf("gen s3 repository encryption: password: %w", err)
	}
	r.Encryption.PasswordEnc = enc

	return nil
}

// initCredential initializes the AWS credentials for the S3Repository.
// It sets the AWS access key ID and secret access key as environment variables.
func (r *S3Repository) initCredential() {
	os.Setenv(awsAccessKeyIdEnv, r.AccessKeyId)
	os.Setenv(awsSecretAccessKeyEnv, r.SecretAccessKey)
}
