package secret

import (
	"fmt"
	"os"

	"github.com/liuminhaw/wrestic-brw/utils/encryptor"
	"github.com/spf13/cobra"
)

// decryptCmd represents the decrypt command
var decryptCmd = &cobra.Command{
	Use:   "decrypt ENC_TEXT [KEY]",
	Short: "Decrypt given text with key in .env file or given key",
	Args:  cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {
		var encKey [32]byte
		var err error
		if len(args) == 2 {
			encKey, err = toEncKey(args[1])
			if err != nil {
				fmt.Printf("Failed to convert given key: %v\n", err)
				os.Exit(1)
			}
		} else {
			encKey, err = loadEncKey()
			if err != nil {
				fmt.Printf("Failed to load enc key: %v\n", err)
				os.Exit(1)
			}
		}

		encText := args[0]
		decText, err := encryptor.Decrypt(encText, encKey)
		if err != nil {
			fmt.Printf("Failed to decrypt: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Decrypted text: %s\n", string(decText))
	},
}

func init() {
	SecretCmd.AddCommand(decryptCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// decryptCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// decryptCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
