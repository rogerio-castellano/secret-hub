package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/rogerio-castellano/secret-hub/internal/crypto"
	"github.com/rogerio-castellano/secret-hub/internal/storage"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	getSecretName string
)

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Retrieve and decrypt a secret by name",
	Long: `Retrieve and decrypt a secret by name using AES-256-GCM.

This command retrieves a secret from the designated storage backend using its name, decrypts it with the provided AES-256-GCM key, and prints the plaintext to standard output. You must specify the name of the secret and a valid decryption key, which can be loaded either from a .secret-hub.yaml configuration file in $HOME or from a standalone 32-byte key file.

Optionally, you may specify a custom path to the secret store file. If no path is given, the default storage backend is used.

Usage 
	get --name <secret_name> [--key <keyfile>] [--storage <filepath>]

Examples:
# Retrieve the secret named 'db_password' from local storage file 'secret-store.json' using the decryption key from 'key.bin'
  secret-hub get --name db_password --key key.bin --storage secret-store.json

# Retrieve the secret named 'db_password' using default configuration and key settings
  secret-hub get --name db_password
`,

	RunE: func(cmd *cobra.Command, args []string) error {
		keyPath := getKey("get")
		storagePath := getStorage("get")
		key, err := crypto.LoadKeyFromFile(keyPath)
		if err != nil {
			return fmt.Errorf("failed to load key: %w", err)
		}

		store := storage.NewFileStore(storagePath)
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

	getCmd.Flags().StringVarP(&getSecretName, "name", "n", "", "Name of the secret (required)")
	getCmd.Flags().StringP("key", "k", "", "Decryption key path (required unless specified in config).")
	getCmd.Flags().StringP("storage", "s", "", "Path to the secret store file")
	if err := viper.BindPFlag("get.key", getCmd.Flags().Lookup("key")); err != nil {
		log.Fatalf("Failed to bind config key: %v", err)
	}
	if err := viper.BindPFlag("get.storage", getCmd.Flags().Lookup("storage")); err != nil {
		log.Fatalf("Failed to bind config key: %v", err)
	}

	if err := getCmd.MarkFlagRequired("name"); err != nil {
		log.Fatalf("Failed to mark flag required: %v", err)
	}
}
