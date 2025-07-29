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
	Long: `The "store" command allows you to securely encrypt a secret value using a provided key
and store it under a specified name. The secret is encrypted with the key loaded from a file,
and then saved to the configured storage backend. If a secret with the same name already exists,
you can use the force flag to overwrite it. This command ensures that sensitive information is
never stored in plaintext, providing an additional layer of security for secret management.

Example:
  secret-hub store --key mykey.bin --name db_password --value "p@ssw0rd"
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
	viper.BindPFlag("store.key", storeCmd.Flags().Lookup("key"))
	viper.BindPFlag("store.storage", storeCmd.Flags().Lookup("storage"))

	storeCmd.MarkFlagRequired("name")
	storeCmd.MarkFlagRequired("value")
}
