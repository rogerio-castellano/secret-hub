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
	Long: `The get command retrieves a secret from the specified store by its name,
decrypts it using the provided key, and prints the plaintext value to stdout.

You must provide the name of the secret, the path to the decryption key, and optionally
the path to the secret store file.

Example:
  secret-hub get --key mykey.bin --name db_password`,
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
