/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/rogerio-castellano/secret-hub/internal/crypto"
	"github.com/rogerio-castellano/secret-hub/internal/storage"
	"github.com/spf13/cobra"
)

var (
	getSecretName string
	getKeyPath    string
	getStorePath  string
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Retrieve and decrypt a secret by name",
	Long: `The get command retrieves a secret from the specified store by its name,
decrypts it using the provided key, and prints the plaintext value to stdout.

You must provide the name of the secret, the path to the decryption key, and optionally
the path to the secret store file.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		key, err := crypto.LoadKeyFromFile(getKeyPath)
		if err != nil {
			return fmt.Errorf("failed to load key: %w", err)
		}

		store := storage.NewFileStore(getStorePath)
		secret, err := store.Get(getSecretName)
		if err != nil {
			return fmt.Errorf("failed to retrieve secret: %w", err)
		}

		plaintext, err := crypto.Decrypt(key, secret.Data)
		if err != nil {
			return fmt.Errorf("decryption failed: %w", err)
		}

		fmt.Fprintln(os.Stdout, string(plaintext))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(getCmd)

	getCmd.Flags().StringVar(&getSecretName, "name", "", "Name of the secret")
	getCmd.Flags().StringVar(&getKeyPath, "key", "", "Path to decryption key")
	getCmd.Flags().StringVar(&getStorePath, "store", "secrets.json", "Path to the secret store file")

	getCmd.MarkFlagRequired("name")
	getCmd.MarkFlagRequired("key")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
