package cmd

import (
	"fmt"
	"log"

	"github.com/rogerio-castellano/secret-hub/internal/storage"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	deleteName string
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a stored secret by name",
	Long: `Delete a secret from the secret store by specifying its name.

This command removes the secret with the given name from the specified store file.
Examples:
  secret-hub delete --name db_password
  secret-hub delete --name mysecret --store secrets.json
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		storagePath := getStorage("delete")
		store := storage.NewFileStore(storagePath)
		if err := store.Delete(deleteName); err != nil {
			return fmt.Errorf("failed to delete secret '%s': %w", deleteName, err)
		}
		fmt.Printf("🗑️  Secret '%s' deleted.\n", deleteName)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)

	deleteCmd.Flags().StringVarP(&deleteName, "name", "n", "", "Name of the secret to delete")
	deleteCmd.Flags().StringP("storage", "s", "", "Path to the secret store file")
	if err := viper.BindPFlag("delete.storage", deleteCmd.Flags().Lookup("storage")); err != nil {
		log.Fatalf("Failed to bind config key: %v", err)
	}

	if err := deleteCmd.MarkFlagRequired("name"); err != nil {
		log.Fatalf("Failed to mark flag required: %v", err)
	}
}
