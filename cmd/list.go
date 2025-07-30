package cmd

import (
	"fmt"
	"log"

	"github.com/rogerio-castellano/secret-hub/internal/storage"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all stored secret names",
	Long: `List all stored secret names without revealing their values.

This command reads the configured secret store file and displays the names of all secrets currently saved. It helps you identify which secrets exist in the store without exposing their content. The storage path can be explicitly defined using a flag or sourced from the configuration file .secret-hub.yaml in $HOME. If no secrets are present, a friendly message is displayed.

Usage 
	list [--storage <filepath>]

Example:
  secret-hub list --storage secret-store.json
  secret-hub list
  `,

	RunE: func(cmd *cobra.Command, args []string) error {
		storagePath := getStorage("list")
		store := storage.NewFileStore(storagePath)

		names, err := store.ListNames()
		if err != nil {
			return fmt.Errorf("failed to list secrets: %w", err)
		}
		if len(names) == 0 {
			fmt.Println("No secrets found.")
			return nil
		}
		fmt.Println("Stored secrets:")
		for _, name := range names {
			fmt.Println(" -", name)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	listCmd.Flags().StringP("storage", "s", "", "Path to the secret store file")
	if err := viper.BindPFlag("list.storage", listCmd.Flags().Lookup("storage")); err != nil {
		log.Fatalf("Failed to bind config key: %v", err)
	}
}
