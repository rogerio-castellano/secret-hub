package cmd

import (
	"fmt"

	"github.com/rogerio-castellano/secret-hub/internal/storage"
	"github.com/spf13/cobra"
)

var listStorePath string

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all stored secret names",
	Long: `List displays the names of all secrets currently stored in the secret store file. 
	Use this command to view which secrets are available without revealing their values.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		store := storage.NewFileStore(listStorePath)

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

	listCmd.Flags().StringVar(&listStorePath, "store", "secrets.json", "Path to the secret store file")
}
