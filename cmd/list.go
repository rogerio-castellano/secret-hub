package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/rogerio-castellano/secret-hub/internal/storage"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	listOutputJSON   bool
	listOutputPretty bool
	listFilterPrefix string
	listQuiet        bool
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all stored secret names",
	Long: `Use the secret-hub list command to enumerate all secrets stored in your configured secret store, without revealing any sensitive values. This is useful for verifying presence, auditing key names, or inspecting store structure.

By default, the command reads from the secret store path specified via a flag or sourced from the .secret-hub.yaml configuration file in $HOME. It displays each stored secret’s name clearly, even when storage is nested or namespaced.

You can tailor the output using:
--json for machine-readable structured listing
--pretty for human-friendly formatted listing (Note: --json and --pretty are mutually exclusive; only one may be used per invocation.)
--filter=<prefix> to display only secrets whose names begin with the given prefix (e.g., --filter=prod_)
--quiet to suppress all output Ideal for:
	- Script automation with exit-code-only checks 
		if ./secret-hub list --filter=prod_ --quiet; then echo "Production secrets are present" fi
	- CI/CD integration where --json or --pretty output is redirected or parsed elsewhere 
		secret-hub list --json --quiet > secrets.json
	- Silent probing or counting operations 
		count=$(./secret-hub list --filter=prod_ | wc -l) [ "$count" -gt 0 ] && echo "✅ Found $count production secrets"
		OR ./secret-hub list --filter=prod_ --quiet && echo "✅ Prod secrets exist"

If no secrets are currently saved, a friendly message is printed—unless suppressed with --quiet.

This command is especially helpful for reviewing what secrets exist, integrating visibility into build pipelines, or validating import results across secure configuration workflows.

Usage 
	list [--storage <filepath>] [--json] [--pretty]

Example:
  secret-hub list --storage secret-store.json
  secret-hub list
  `,

	RunE: func(cmd *cobra.Command, args []string) error {
		if listOutputJSON && listOutputPretty {
			return fmt.Errorf("flags --json and --pretty are mutually exclusive; please choose only one")
		}

		storagePath := getStorage("list")
		store := storage.NewFileStore(storagePath)

		names, err := store.List()
		if err != nil {
			return fmt.Errorf("failed to list secrets: %w", err)
		}

		// Apply prefix filtering if set
		if listFilterPrefix != "" {
			filtered := make([]string, 0, len(names))
			for _, name := range names {
				if strings.HasPrefix(name, listFilterPrefix) {
					filtered = append(filtered, name)
				}
			}
			names = filtered
		}

		if listQuiet {
			return nil
		}

		if listOutputJSON {
			out, err := json.MarshalIndent(names, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to marshal JSON: %w", err)
			}
			fmt.Println(string(out))
			return nil
		}

		if len(names) == 0 {
			fmt.Println("No secrets found.")
			return nil
		}

		if listOutputPretty || (!listOutputJSON && !listOutputPretty) {
			fmt.Println("🔐 Secrets in store:")
			fmt.Println("--------------------")
			for _, name := range names {
				fmt.Println("•", name)
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	listCmd.Flags().StringP("storage", "s", "", "Path to the secret store file")
	listCmd.Flags().BoolVar(&listOutputJSON, "json", false, "Print output as JSON array")
	listCmd.Flags().BoolVar(&listOutputPretty, "pretty", false, "Print formatted output (default)")
	listCmd.Flags().BoolVar(&listQuiet, "quiet", false, "Suppress all output of list (useful for script automation and CI/CD)")
	listCmd.Flags().StringVar(&listFilterPrefix, "filter", "", "Optional prefix filter (e.g., 'prod_')")

	if err := viper.BindPFlag("list.storage", listCmd.Flags().Lookup("storage")); err != nil {
		log.Fatalf("Failed to bind config key: %v", err)
	}
}
