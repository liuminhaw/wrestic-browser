/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package secret

import (
	"fmt"
	"os"

	"github.com/liuminhaw/wrestic-brw/utils/encryptor"
	"github.com/spf13/cobra"
)

// encryptCmd represents the encrypt command
var encryptCmd = &cobra.Command{
	Use:   "encrypt TEXT [KEY]",
	Short: "Encrypt given text with key in .env file or given key",
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

		message := args[0]
		encText, err := encryptor.Encrypt([]byte(message), encKey)
		if err != nil {
			fmt.Printf("Failed to encrypt: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Encrypted text: %s\n", encText)
	},
}

func init() {
	SecretCmd.AddCommand(encryptCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// encryptCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// encryptCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
