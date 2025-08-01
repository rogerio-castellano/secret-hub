package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/rogerio-castellano/secret-hub/internal/crypto"
	"github.com/rogerio-castellano/secret-hub/internal/storage"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

var (
	exportFormat        string
	exportSummaryFormat string
	exportQuiet         bool
)

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export decrypted secrets in .env, JSON, or YAML format",
	Long: `Export all stored secrets in decrypted form using a selected output format.

Supports .env, JSON, and YAML formats to accommodate different integration needs.
This is useful for injecting secrets into local development environments,
generating configuration files, or importing into other secret management systems.

Additional flags:
--summary  Provides a concise overview of exported secrets without showing values.
--quiet    Suppresses all non-essential output for seamless scripting and automation.

Usage:
  export --format <env|json|yaml> [--key <keyfile>] [--storage <filepath>] [--summary <json|text>] [--quiet]

Examples:
  secret-hub export --format env --key key.bin
  secret-hub export --format json --key key.bin > exported-secrets.json
  secret-hub export --format yaml --key key.bin --storage secret-store.json
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		keyPath := getKey("export")
		storagePath := getStorage("export")
		key, err := crypto.LoadKeyFromFile(keyPath)
		if err != nil {
			return fmt.Errorf("failed to load key: %w", err)
		}

		store := storage.NewFileStore(storagePath)
		names, err := store.ListNames()
		if err != nil {
			return fmt.Errorf("failed to list secret names: %w", err)
		}

		output := make(map[string]string)

		for _, name := range names {
			encSecret, err := store.Get(name)
			if err != nil {
				return fmt.Errorf("failed to get secret %s: %w", name, err)
			}
			plain, err := crypto.Decrypt(key, encSecret.Data)
			if err != nil {
				return fmt.Errorf("failed to decrypt %s: %w", name, err)
			}
			output[name] = string(plain)
		}
		count := len(output)

		switch exportFormat {
		case "env":
			for k, v := range output {
				fmt.Printf("%s=%s\n", k, v)
			}
		case "json":
			out, err := json.MarshalIndent(output, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to marshal JSON: %w", err)
			}
			os.Stdout.Write(out)
		case "yaml":
			out, err := yaml.Marshal(output)
			if err != nil {
				return fmt.Errorf("failed to marshal YAML: %w", err)
			}
			os.Stdout.Write(out)
		default:
			return fmt.Errorf("unsupported format: %s (use env, json, yaml)", exportFormat)
		}

		if exportQuiet {
			return nil
		}

		summary := struct {
			Store   string `json:"store"`
			Format  string `json:"format"`
			Secrets int    `json:"secrets_exported"`
		}{
			Store:   storagePath,
			Format:  exportFormat,
			Secrets: count,
		}

		switch exportSummaryFormat {
		case "json":
			out, err := json.MarshalIndent(summary, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to encode summary: %w", err)
			}
			fmt.Println("\n📦 Export Summary")
			fmt.Println(string(out))
		default:
			fmt.Println("\n📦 Export Summary")
			fmt.Println("-----------------")
			fmt.Printf("📁 Store       : %s\n", summary.Store)
			fmt.Printf("📤 Format      : %s\n", summary.Format)
			fmt.Printf("🔐 Secrets     : %d\n", summary.Secrets)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(exportCmd)

	exportCmd.Flags().StringVarP(&exportFormat, "format", "f", "env", "Export format: env, json, yaml")
	exportCmd.Flags().StringP("key", "k", "", "Path to decryption key file")
	exportCmd.Flags().StringP("storage", "s", "", "Path to secret storage file")
	exportCmd.Flags().StringVar(&exportSummaryFormat, "summary", "text", "Summary output format: text or json")
	exportCmd.Flags().BoolVar(&exportQuiet, "quiet", false, "Suppress export summary output")

	if err := viper.BindPFlag("export.key", exportCmd.Flags().Lookup("key")); err != nil {
		log.Fatalf("Failed to bind config key: %v", err)
	}
	if err := viper.BindPFlag("export.storage", exportCmd.Flags().Lookup("storage")); err != nil {
		log.Fatalf("Failed to bind config key: %v", err)
	}
}
