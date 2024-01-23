package restic

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"time"
)

type LocalRepository struct {
	Password    string
	Destination string
}

func (r *LocalRepository) Connect() error {
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
