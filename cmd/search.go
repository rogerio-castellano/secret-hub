package cmd

import (
	"fmt"
	"log"
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
	Long: `Search stored secret names by substring match.

This command scans all secret names stored in the configured secret file and returns those that include the specified substring, matched in a case-insensitive manner. It allows quick discovery of secrets based on partial naming patterns.

By default, the secret store is resolved from the .secret-hub.yaml configuration file in $HOME. You can also explicitly define the storage path with the --storage flag.

If no matching names are found, a friendly notice will be displayed. This command is useful for locating secrets without decrypting them or exposing their contents.

Usage 
	search --query <substring> [--storage <filepath>]

Example:
  secret-hub search --query db --storage secret-store.json
  secret-hub search --query db
`,
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
	if err := viper.BindPFlag("search.storage", searchCmd.Flags().Lookup("storage")); err != nil {
		log.Fatalf("Failed to bind config key: %v", err)
	}

	if err := searchCmd.MarkFlagRequired("query"); err != nil {
		log.Fatalf("Failed to bind config key: %v", err)
	}
}
