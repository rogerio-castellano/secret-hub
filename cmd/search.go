package cmd

import (
	"fmt"
	"strings"

	"github.com/rogerio-castellano/secret-hub/internal/storage"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	searchQuery string
)

var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search for stored secret names by substring",
	Long: `Search stored secrets for names that match a given substring.

Example:
  secret-hub search --query TOKEN
  secret-hub search --query email --store my-secrets.json`,
	RunE: func(cmd *cobra.Command, args []string) error {
		storagePath := getStorage("search")
		store := storage.NewFileStore(storagePath)

		names, err := store.ListNames()
		if err != nil {
			return fmt.Errorf("failed to list secrets: %w", err)
		}

		var matches []string
		for _, name := range names {
			if strings.Contains(strings.ToLower(name), strings.ToLower(searchQuery)) {
				matches = append(matches, name)
			}
		}

		if len(matches) == 0 {
			fmt.Println("🔍 No matching secrets found.")
			return nil
		}

		fmt.Println("🔍 Matching secrets:")
		for _, name := range matches {
			fmt.Println(" -", name)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)

	searchCmd.Flags().StringVarP(&searchQuery, "query", "q", "", "Substring to search for (case-insensitive)")
	searchCmd.Flags().StringP("storage", "s", "", "Path to the secret store file")
	viper.BindPFlag("search.storage", searchCmd.Flags().Lookup("storage"))

	searchCmd.MarkFlagRequired("query")
}
