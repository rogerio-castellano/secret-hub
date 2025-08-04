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
	Long: `Delete a stored secret by name from the configured secret store.

This command removes a secret from the storage backend by specifying its unique name. You must provide the name of the secret to delete. By default, the secret store path is loaded from the .secret-hub.yaml configuration file in $HOME. If this file is unavailable, a custom storage path can be provided using the --storage flag.

If the specified secret does not exist, an error message will be displayed. Successfully deleted secrets will be acknowledged with a confirmation output.

Usage 
	delete --name <secret_name> [--storage <filepath>]

Examples:
# Delete the secret named 'mysecret' from the 'secret-store.json' storage file
  secret-hub delete --name mysecret --storage secret-store.json

# Delete the secret named 'db_password' using the default storage location
  secret-hub delete --name db_password

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
