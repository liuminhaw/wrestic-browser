package restic

import (
	"context"
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
	Password        string
	Destination     string
	AccessKeyId     string
	SecretAccessKey string
}

func (r *S3Repository) Connect() error {
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

func (r *S3Repository) initCredential() {
	os.Setenv(awsAccessKeyIdEnv, r.AccessKeyId)
	os.Setenv(awsSecretAccessKeyEnv, r.SecretAccessKey)
}
