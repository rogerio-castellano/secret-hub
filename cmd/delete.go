/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/rogerio-castellano/secret-hub/internal/storage"
	"github.com/spf13/cobra"
)

var (
	deleteName      string
	deleteStorePath string
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a stored secret by name",
	Long: `Delete a secret from the secret store by specifying its name.

This command removes the secret with the given name from the specified store file.
Example:
  secret-hub delete --name mysecret --store secrets.json
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		store := storage.NewFileStore("secrets.json")

		if err := store.Delete(deleteName); err != nil {
			return fmt.Errorf("failed to delete secret '%s': %w", deleteName, err)
		}
		fmt.Printf("🗑️  Secret '%s' deleted.\n", deleteName)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)

	deleteCmd.Flags().StringVar(&deleteName, "name", "", "Name of the secret to delete")
	deleteCmd.Flags().StringVar(&deleteStorePath, "store", "secrets.json", "Path to the secret store file")

	deleteCmd.MarkFlagRequired("name")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deleteCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deleteCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
