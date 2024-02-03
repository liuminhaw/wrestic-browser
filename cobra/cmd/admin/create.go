/*
Copyright © 2023 Min-Haw, Liu

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package admin

import (
	"fmt"
	"os"
	"syscall"

	"github.com/liuminhaw/wrestic-brw/cobra/cmd/password"
	"github.com/liuminhaw/wrestic-brw/models"
	"github.com/liuminhaw/wrestic-brw/utils/dotenv"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create admin account",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := loadEnvConfig()
		if err != nil {
			fmt.Printf("Failed to load dotenv: %v\n", err)
			os.Exit(1)
		}

		db, err := models.Open(cfg.PSQL)
		if err != nil {
			fmt.Printf("Failed to create db connection: %v\n", err)
			os.Exit(1)
		}
		defer db.Close()

		// Ask for admin username and password
		var username string
		fmt.Print("Enter admin username: ")
		fmt.Scanln(&username)

		fmt.Print("Enter admin password (will be hidden): ")
		bytePassword, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			fmt.Printf("Failed to read password: %v\n", err)
			os.Exit(1)
		}
		fmt.Print("\nEnter admin password again (will be hidden): ")
		bytePasswordConfirm, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			fmt.Printf("Failed to read password: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("\n")
		if string(bytePassword) != string(bytePasswordConfirm) {
			fmt.Println("Input password does not match")
			os.Exit(1)
		}

		passwordHash, err := password.Hash(string(bytePassword))
		if err != nil {
			fmt.Printf("Failed to hash password: %v\n", err)
			os.Exit(1)
		}

		// Insert admin to Database
		_, err = db.Exec(`
			INSERT INTO users (username, password_hash, role_id)
			VALUES ($1, $2, (SELECT id FROM roles WHERE role = 'admin'));
		`, username, passwordHash)
		if err != nil {
			fmt.Printf("Failed to create admin user: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Admin user %s created\n", username)
	},
}

func init() {
	AdminCmd.AddCommand(createCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

type config struct {
	PSQL models.PostgresConfig
}

// loadEnvConfig loads config settings from .env file
func loadEnvConfig() (config, error) {
	var cfg config
	err := dotenv.LoadDotEnv()
	if err != nil {
		return cfg, fmt.Errorf("load env config: %w", err)
	}

	// Read PostgreSQL values from env variables
	cfg.PSQL.Host = os.Getenv("DB_HOST")
	cfg.PSQL.Port = os.Getenv("DB_PORT")
	cfg.PSQL.User = os.Getenv("DB_USER")
	cfg.PSQL.Password = os.Getenv("DB_PASSWORD")
	cfg.PSQL.Database = os.Getenv("DB_DATABASE")
	cfg.PSQL.SSLMode = os.Getenv("DB_SSLMODE")

	return cfg, nil
}
