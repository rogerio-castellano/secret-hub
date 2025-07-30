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
	importFilePath string
	importFormat   string
	importForce    bool
)

var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Import secrets from .env, JSON, or YAML and store them encrypted",
	Long: `Securely import secrets from a structured file and store them in encrypted format using the secret-hub import command. This operation helps you onboard existing environment variables or configuration values into the secure storage backend—whether from .env, JSON, or YAML files.

Each value in the input file is encrypted using the key specified via the --key flag. You can also specify the destination file for the encrypted secrets using --storage, or allow the tool to fallback to the default configured backend. Existing secrets can be forcibly overwritten using the --force flag.

The import process supports flexible file formats and provides error handling for parsing, encryption, and storage operations. Ideal for teams migrating legacy secrets or initializing fresh environments with sensitive configs.

Usage:
  secret-hub import --file <path> --format <env|json|yaml> --key <keyfile> [--storage <filepath>] [--force]

Examples:
  secret-hub import --file secrets.env --format env --key test-key.bin --storage imported-secrets.json 
  secret-hub import --file secrets.json --format json --key test-key.bin --force
  secret-hub import --file secrets.yaml --format yaml
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

		key, err := crypto.LoadKeyFromFile(keyPath)
		if err != nil {
			return fmt.Errorf("failed to load key: %w", err)
		}

		store := storage.NewFileStore(storagePath)
		imported := 0

		for name, value := range data {
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
