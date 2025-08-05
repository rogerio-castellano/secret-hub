package cmd

import (
	"fmt"
	"log"

	"github.com/rogerio-castellano/secret-hub/internal/crypto"
	"github.com/rogerio-castellano/secret-hub/internal/storage"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	secretName  string
	secretValue string
	storeKey    string
	forceStore  bool
	storePath   string
)

var storeCmd = &cobra.Command{
	Use:   "store",
	Short: "Encrypt and store a secret by name",
	Long: `Encrypt and store a secret by name.

This command securely encrypts a secret value with the AES-256-GCM algorithm and stores it under a unique name. 
You must provide both the secret value and its corresponding name, as well as an encryption key sourced either 
from a .secret-hub.yaml config file in $HOME or from a separate key file containing a valid 32-byte key. 

The encrypted secret is saved to the specified storage backend (default: secrets.json). 
If a secret with the same name already exists, use the --force flag to overwrite it.

Usage:
  store --name <secret_name> --value <secret_value> [--key <keyfile>] [--storage <filepath>] [--force]

Examples:
# Store the secret 'db_password' with the value "p@ssw0rd" in 'secret-store.json', using 'key.bin' to encrypt
  secret-hub store --key key.bin --storage=secret-store.json --name db_password --value "p@ssw0rd"

# Store the secret 'db_password' with the value "p@ssw0rd" using default encryption key and storage settings
  secret-hub store --name db_password --value "p@ssw0rd"

# Overwrite existing 'db_password' secret with the new value "p@ssw0rd" in 'secret-store.json', using 'key.bin' and --force to overwrite if a secret with the same name already exists
  secret-hub store --key key.bin --storage=secret-store.json --name db_password --value "p@ssw0rd" --force

`,
	RunE: func(cmd *cobra.Command, args []string) error {
		storeKey = getKey("store")
		storePath = getStorage("store")
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

	storeCmd.Flags().StringVarP(&secretName, "name", "n", "", "Name of the secret to store (required)")
	storeCmd.Flags().StringVarP(&secretValue, "value", "l", "", "Value of the secret to store (required)")
	storeCmd.Flags().StringVarP(&storeKey, "key", "k", "", "Encryption key path (required unless specified in config).")
	storeCmd.Flags().StringVarP(&storePath, "storage", "s", "", "Path to the storage file (default to secrets.json)")
	storeCmd.Flags().BoolVarP(&forceStore, "force", "f", false, "Force overwrite existing secret")
	if err := viper.BindPFlag("store.key", storeCmd.Flags().Lookup("key")); err != nil {
		log.Fatalf("Failed to bind config key: %v", err)
	}
	if err := viper.BindPFlag("store.storage", storeCmd.Flags().Lookup("storage")); err != nil {
		log.Fatalf("Failed to bind config key: %v", err)
	}

	if err := storeCmd.MarkFlagRequired("name"); err != nil {
		log.Fatalf("Failed to mark flag required: %v", err)
	}
	if err := storeCmd.MarkFlagRequired("value"); err != nil {
		log.Fatalf("Failed to mark flag required: %v", err)
	}
}
