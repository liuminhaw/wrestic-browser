package secret

import (
	"fmt"
	"os"

	"github.com/liuminhaw/wrestic-brw/rand"
	"github.com/spf13/cobra"
)

// genKeyCmd represents the genKey command
var genKeyCmd = &cobra.Command{
	Use:   "genKey",
	Short: "Generate a new key for encryption and decryption via secretbox",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		key, err := rand.String(32)
		if err != nil {
			fmt.Printf("Error generating key: %s\n", err)
			os.Exit(1)
		}
		fmt.Printf("Generated key: %s\n", key)
	},
}

func init() {
	SecretCmd.AddCommand(genKeyCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// genKeyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// genKeyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
