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
	importFilePath     string
	importFormat       string
	importForce        bool
	importDryRun       bool
	importSkipExisting bool
	importOnlyNames    []string
)

var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Import secrets from .env, JSON, or YAML and store them encrypted",
	Long: `Securely import secrets from a structured file and store them in encrypted format using the secret-hub import command. This operation helps onboard existing environment variables or configuration values into the secure storage backend—whether from .env, JSON, or YAML files.

Each value in the input file is encrypted using the key specified via the --key flag. You can define the output destination using --storage, or allow the tool to fall back to the default configured backend. To prevent overwriting existing secrets, use the --skip-existing flag—even without --force, this ensures only new secrets are imported. If overwriting is intended, --force makes that explicit.

To preview secret names and values without encrypting or saving them, use the --dry-run flag. This mode simulates the import process, providing visibility into the parsed data and helping validate input before committing changes.

For targeted imports, the --only flag allows you to selectively process specific secrets from a larger file—useful when reusing files across environments or importing just a subset of values. You can pass a comma-separated list of keys to limit the operation to just those entries.

The import process supports flexible file formats and includes robust error handling for parsing, encryption, and storage operations. Ideal for teams migrating legacy secrets, auditing secret inventories, or initializing secure development environments.

Usage:
  secret-hub import --file <path> --format <env|json|yaml> --key <keyfile> 
                  [--storage <filepath>] [--force] [--skip-existing] [--dry-run] [--only <keys>]

Examples:
# Import from a .env file into custom storage
  secret-hub import --file secrets.env --format env --key test-key.bin --storage imported-secrets.json 

# Import from JSON and overwrite existing secrets
  secret-hub import --file secrets.json --format json --key test-key.bin --force

# Import from YAML and skip secrets that already exist
  secret-hub import --file secrets.yaml --format yaml --skip-existing

# Dry-run preview without encrypting or saving
  secret-hub import --file secrets.json --format json --key test-key.bin --dry-run

# import only specific secrets (e.g., API_KEY and DB_PASSWORD):
  secret-hub import config.env --key superkey123 --only API_KEY,DB_PASSWORD  
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		keyPath := getKey("import")
		storagePath := getStorage("import")

		content, err := os.ReadFile(importFilePath)
		if err != nil {
			return fmt.Errorf("failed to read file: %w", err)
		}

		var data map[string]string

		switch importFormat {
		case "env":
			data = parseEnv(string(content))
		case "json":
			err := json.Unmarshal(content, &data)
			if err != nil {
				return fmt.Errorf("invalid JSON: %w", err)
			}
		case "yaml":
			err := yaml.Unmarshal(content, &data)
			if err != nil {
				return fmt.Errorf("invalid YAML: %w", err)
			}
		default:
			return fmt.Errorf("unsupported format: %s (use env, json, yaml)", importFormat)
		}

		if len(data) == 0 {
			return fmt.Errorf("no secrets found in file")
		}

		// Filter keys if --only is set
		if len(importOnlyNames) > 0 {
			filtered := make(map[string]string)
			for _, name := range importOnlyNames {
				if val, ok := data[name]; ok {
					filtered[name] = val
				} else {
					fmt.Printf("⚠️  Warning: --only specified key '%s' not found in input file\n", name)
				}
			}
			data = filtered
		}

		key, err := crypto.LoadKeyFromFile(keyPath)
		if err != nil {
			return fmt.Errorf("failed to load key: %w", err)
		}

		store := storage.NewFileStore(storagePath)
		imported := 0

		for name, value := range data {
			// Check existence
			if importSkipExisting && !importForce {
				existing, _ := store.Get(name)
				if existing != nil {
					fmt.Printf("⚠️  Skipping existing secret: %s\n", name)
					continue
				}
			}

			if importDryRun {
				fmt.Printf("🔍 [dry-run] Would import: %s = %s\n", name, value)
				continue
			}

			enc, err := crypto.Encrypt(key, []byte(value))
			if err != nil {
				return fmt.Errorf("failed to encrypt secret '%s': %w", name, err)
			}

			s := storage.EncryptedSecret{
				Name: name,
				Data: enc,
			}

			if err := store.Save(s, importForce); err != nil {
				return fmt.Errorf("failed to store '%s': %w", name, err)
			}

			imported++
		}

		if importDryRun {
			log.Printf("✅ Dry run complete. %d secrets would be imported from %s", len(data), importFilePath)
			return nil
		}

		log.Printf("✅ Imported %d secrets from %s", imported, importFilePath)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(importCmd)

	importCmd.Flags().StringVar(&importFilePath, "file", "", "Path to input file (required)")
	importCmd.Flags().StringVar(&importFormat, "format", "", "Input format: env, json, yaml (required)")
	importCmd.Flags().StringP("key", "k", "", "Path to encryption key file (required)")
	importCmd.Flags().StringP("storage", "s", "", "Path to secret store file")
	importCmd.Flags().BoolVar(&importForce, "force", false, "Overwrite existing secrets")
	importCmd.Flags().BoolVar(&importDryRun, "dry-run", false, "Show what would be imported without writing")
	importCmd.Flags().BoolVar(&importSkipExisting, "skip-existing", false, "Skip importing secrets that already exist")
	importCmd.Flags().StringSliceVar(&importOnlyNames, "only", nil, "Comma-separated list of secret names to import")

	if err := viper.BindPFlag("import.key", importCmd.Flags().Lookup("key")); err != nil {
		log.Fatalf("Failed to bind config key: %v", err)
	}
	if err := viper.BindPFlag("import.storage", importCmd.Flags().Lookup("storage")); err != nil {
		log.Fatalf("Failed to bind config key: %v", err)
	}

	if err := importCmd.MarkFlagRequired("file"); err != nil {
		log.Fatalf("Failed to mark flag required: %v", err)
	}
	if err := importCmd.MarkFlagRequired("format"); err != nil {
		log.Fatalf("Failed to mark flag required: %v", err)
	}
}

func parseEnv(content string) map[string]string {
	result := make(map[string]string)
	lines := strings.SplitSeq(content, "\n")
	for line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		result[parts[0]] = parts[1]
	}
	return result
}
