package restic

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"time"
)

type SftpRepository struct {
	Password    string
	Destination string
	User        string
	Host        string
	Pem         string
}

func (r *SftpRepository) Connect() error {
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
