package secret

import (
	"encoding/base64"
	"fmt"
	"os"

	"github.com/liuminhaw/wrestic-brw/utils/dotenv"
	"github.com/spf13/cobra"
)

// secretCmd represents the secret command
var SecretCmd = &cobra.Command{
	Use:   "secret",
	Short: "Secret related command like encryption and decryption",
	Long:  ``,
	// Run: func(cmd *cobra.Command, args []string) {
	// 	fmt.Println("secret called")
	// },
}

func init() {
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// secretCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// secretCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// loadEncKey loads the encryption key from the environment variables.
// It returns the encryption key as a [32]byte array and an error if any.
func loadEncKey() ([32]byte, error) {
	err := dotenv.LoadDotEnv()
	if err != nil {
		return [32]byte{}, fmt.Errorf("load enc key: %w", err)
	}

	encKey := os.Getenv("ENC_KEY")
	if encKey == "" {
		return [32]byte{}, fmt.Errorf("load enc key: empty enc key")
	}
	encKeyBytes, err := base64.URLEncoding.DecodeString(encKey)
	if err != nil {
		return [32]byte{}, fmt.Errorf("load enc key: decode enc key: %w", err)
	}

	return [32]byte(encKeyBytes), nil
}

// toEncKey converts a base64-encoded key string to a [32]byte key.
// It decodes the key using base64.URLEncoding and returns the key as a [32]byte.
// If there is an error decoding the key, it returns an empty [32]byte and an error.
func toEncKey(key string) ([32]byte, error) {
	encKeyBytes, err := base64.URLEncoding.DecodeString(key)
	if err != nil {
		return [32]byte{}, fmt.Errorf("toEncKey: decode key: %w", err)
	}

	return [32]byte(encKeyBytes), nil
}
