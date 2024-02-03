package dotenv

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

// LoadDotEnv loads the environment variables from a .env file located in the same directory as the executable.
// It returns an error if it fails to get the executable path or fails to load the .env file.
func LoadDotEnv() error {
	execPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("load dotenv: failed to get executable: %w", err)
	}
	execDir := filepath.Dir(execPath)

	err = godotenv.Load(filepath.Join(execDir, ".env"))
	if err != nil {
		return fmt.Errorf("load dotenv: failed to load: %w", err)
	}

	return nil
}
