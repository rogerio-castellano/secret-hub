package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

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
	exportOutputPath    string
)

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export decrypted secrets in .env, JSON, or YAML format",
	Long: `Export all stored secrets in decrypted form using a selected output format.

Supports .env, JSON, and YAML formats to accommodate different integration needs—useful for injecting secrets into development environments, generating config files, or migrating secrets into other management systems.

If --format is omitted, the tool will intelligently infer the desired output format from the --output file's extension (e.g. .json, .yaml, .env). This streamlines automation and scripting by reducing boilerplate setup.

Additional flags:

--output <filename> Saves exported secrets to a file instead of printing to stdout. Ideal for CI pipelines and provisioning tasks.
--summary Outputs a brief, non-sensitive overview of what was exported.
--quiet Suppresses all non-essential output for silent scripting workflows.
--format <format> Explicitly defines the format if auto-detection isn't desired.

Usage:
  export --format <env|json|yaml> [--key <keyfile>] [--storage <filepath>] [--summary <json|text>] [--quiet]
  [--output <filename>]

Examples:
  secret-hub export --format env --key key.bin
  secret-hub export --format json --key key.bin > exported-secrets.json
  secret-hub export --format yaml --key key.bin --storage secret-store.json
  secret-hub export --format json --output .env
  secret-hub export --format env --summary --output .env.generated
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

		// Auto-detect export format from output file extension
		if exportFormat == "" && exportOutputPath != "" {
			switch {
			case strings.HasSuffix(exportOutputPath, ".env"):
				exportFormat = "env"
			case strings.HasSuffix(exportOutputPath, ".json"):
				exportFormat = "json"
			case strings.HasSuffix(exportOutputPath, ".yml"), strings.HasSuffix(exportOutputPath, ".yaml"):
				exportFormat = "yaml"
			default:
				return fmt.Errorf("cannot auto-detect format from output file extension. Use --format explicitly")
			}
		}

		var out []byte
		switch exportFormat {
		case "env":
			var sb strings.Builder
			for k, v := range output {
				sb.WriteString(fmt.Sprintf("%s=%s\n", k, v))
			}
			out = []byte(sb.String())
		case "json":
			var err error
			out, err = json.MarshalIndent(output, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to marshal JSON: %w", err)
			}
		case "yaml":
			var err error
			out, err = yaml.Marshal(output)
			if err != nil {
				return fmt.Errorf("failed to marshal YAML: %w", err)
			}
		default:
			return fmt.Errorf("unsupported format: %s (use env, json, yaml)", exportFormat)
		}

		if exportOutputPath != "" {
			err := os.WriteFile(exportOutputPath, out, 0644)
			if err != nil {
				return fmt.Errorf("failed to write output file: %w", err)
			}
			if !exportQuiet {
				fmt.Printf("💾 Secrets written to: %s\n", exportOutputPath)
			}
		} else {
			os.Stdout.Write(out)
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

	exportCmd.Flags().StringVar(&exportFormat, "format", "", "Export format: env, json, yaml (optional if --output is set)")
	exportCmd.Flags().StringP("key", "k", "", "Path to decryption key file")
	exportCmd.Flags().StringP("storage", "s", "", "Path to secret storage file")
	exportCmd.Flags().StringVar(&exportSummaryFormat, "summary", "text", "Summary output format: text or json")
	exportCmd.Flags().BoolVar(&exportQuiet, "quiet", false, "Suppress export summary output")
	exportCmd.Flags().StringVar(&exportOutputPath, "output", "", "Write exported secrets to file instead of stdout")

	if err := viper.BindPFlag("export.key", exportCmd.Flags().Lookup("key")); err != nil {
		log.Fatalf("Failed to bind config key: %v", err)
	}
	if err := viper.BindPFlag("export.storage", exportCmd.Flags().Lookup("storage")); err != nil {
		log.Fatalf("Failed to bind config key: %v", err)
	}
}
