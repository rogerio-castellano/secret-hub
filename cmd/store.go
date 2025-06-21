/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"

	"github.com/rogerio-castellano/secret-hub/internal/crypto"
	"github.com/rogerio-castellano/secret-hub/internal/storage"
	"github.com/spf13/cobra"
)

var (
	secretName  string
	secretValue string
	storeKey    string
	forceStore  bool
	storePath   string
)

// storeCmd represents the store command
// storeCmd represents the "store" command for the CLI application.
//
// Short: Encrypt and store a secret by name.
var storeCmd = &cobra.Command{
	Use:   "store",
	Short: "Encrypt and store a secret by name",
	Long: `The "store" command allows you to securely encrypt a secret value using a provided key
and store it under a specified name. The secret is encrypted with the key loaded from a file,
and then saved to the configured storage backend. If a secret with the same name already exists,
you can use the force flag to overwrite it. This command ensures that sensitive information is
never stored in plaintext, providing an additional layer of security for secret management. `,
	RunE: func(cmd *cobra.Command, args []string) error {
		key, err := crypto.LoadKeyFromFile(storeKey)
		if err != nil {
			return fmt.Errorf("failed to load key: %w", err)
		}
		encrypted, err := crypto.Encrypt(key, []byte(secretValue))
		if err != nil {
			return fmt.Errorf("encryption failed: %w", err)
		}
		store := storage.NewFileStore(storePath)
		sec := storage.EncryptedSecret{
			Name: secretName,
			Data: encrypted,
		}
		if err := store.Save(sec, forceStore); err != nil {
			return fmt.Errorf("failed to store secret: %w", err)
		}
		log.Printf("🔐 Secret stored successfully.")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(storeCmd)

	storeCmd.Flags().StringVar(&secretName, "name", "", "Name of the secret to store (required)")
	storeCmd.Flags().StringVar(&secretValue, "value", "", "Value of the secret to store (required)")
	storeCmd.Flags().StringVar(&storeKey, "key", "", "Path to the encryption key file (required)")
	storeCmd.Flags().StringVar(&storePath, "store", "secrets.json", "Path to the storage file (required)")
	storeCmd.Flags().BoolVar(&forceStore, "force", false, "Force overwrite existing secret")

	storeCmd.MarkFlagRequired("name")
	storeCmd.MarkFlagRequired("value")
	storeCmd.MarkFlagRequired("key")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// storeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// storeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
