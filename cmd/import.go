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
	importFilePath      string
	importFormat        string
	importForce         bool
	importDryRun        bool
	importSkipExisting  bool
	importOnlyNames     []string
	importPrefix        string
	importRenameMap     map[string]string
	importRenameList    []string
	importExclude       []string
	importSummaryFormat string
	importQuiet         bool
)

var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Import secrets from .env, JSON, or YAML and store them encrypted",
	Long: `Securely import secrets from a structured file and store them in encrypted format using the secret-hub import command. This operation helps onboard existing environment variables or configuration values into the secure storage backend—whether from .env, JSON, or YAML files.

Each value in the input file is encrypted using the key specified via the --key flag. You can define the output destination using --storage, or allow the tool to fall back to the default configured backend. To prevent overwriting existing secrets unintentionally, use --skip-existing; for intentional replacements, use --force.

Control over the import scope is provided through the --only flag to include specific secrets and --exclude to skip ones you don't want imported. You can organize secrets using --prefix, remap individual key names using --rename, and preview the import behavior without saving anything via the --dry-run flag.

If --format is not explicitly set, the tool will auto-detect the format from the file extension provided via --file—reducing manual configuration and smoothing CI/CD workflows. Supported extensions include .env, .json, and .yaml.

After import completes, a summary of the operation is shown by default—listing how many secrets were imported, skipped, excluded, renamed, or failed. This can be printed in machine-readable format using --summary=json, or silenced entirely with --quiet for minimal output.

The import command supports flexible formats and provides robust error handling across parsing, encryption, and storage—making it ideal for managing environment configs, migrating shared secret sets, and maintaining clarity across secure systems.

Usage:
  secret-hub import --file <path> --format <env|json|yaml> --key <keyfile> 
                  [--storage <filepath>] [--force] [--skip-existing] [--dry-run] [--only <keys>] [--rename <keys>]
				  [--exclude <keys>] [--summary=json] [--quiet]
				  
Examples:
# Import from a .env file into custom storage
  secret-hub import --file secrets.env --format env --key test-key.bin --storage imported-secrets.json 

# Import from JSON and overwrite existing secrets
  secret-hub import --file secrets.json --format json --key test-key.bin --force

# Import from YAML and skip secrets that already exist
  secret-hub import --file secrets.yaml --format yaml --skip-existing

# Dry-run preview without encrypting or saving
  secret-hub import --file secrets.json --format json --key test-key.bin --dry-run

# Import only specific secrets (e.g., API_KEY and DB_PASSWORD)
  secret-hub import config.env --key superkey123 --only API_KEY,DB_PASSWORD  

# Import prepending dev_ to each key (e.g. DB_PASS becomes dev_DB_PASS)
  secret-hub import --file ./config.env --key masterKey123 --prefix dev_

# Import selectively  with prefixing; only DB_PASS and API_KEY will be imported, with keys remapped to prod_DB_PASS and prod_API_KEY.
  secret-hub import --file ./config.env --key masterKey123 --only DB_PASS,API_KEY --prefix prod_

# Dry run with preview of prefixed output. Preview what will be saved as test_-prefixed keys, allowing validation before committing.
  secret-hub import --file ./secrets.yaml --key masterKey123 --prefix test_ --dry-run

# Import the secrets from secrets.env, renaming DB_PASS to DATABASE_PASSWORD and API_KEY to SERVICE_API_KEY.
  secret-hub import secrets.env --rename DB_PASS=DATABASE_PASSWORD --rename API_KEY=SERVICE_API_KEY

# Import only DB_USER and DB_PASS, and their keys are renamed to USERNAME and PASSWORD
  secret-hub import config.yaml --only DB_USER DB_PASS --rename DB_USER=USERNAME --rename DB_PASS=PASSWORD

# Import the TOKEN secret as prod_ACCESS_TOKEN in storage, utilizing the combined effect of --prefix and --rename.
  secret-hub import app.env --prefix prod_ --rename TOKEN=ACCESS_TOKEN

# Dry run prefixing each key with dev_ (e.g., dev_API_KEY), renaming DATABASE_URL to db_url, and excluding DEBUG_MODE and LOCAL_SECRET—even if they're in the file.
  secret-hub import --file=secrets.env --format env --prefix dev_ --rename DATABASE_URL=db_url --exclude DEBUG_MODE,LOCAL_SECRET --dry-run --file=secrets.env
  
# Import JSON secrets with a summary in machine-readable format:
  secret-hub import --file secrets.json --key team-key --summary=json
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var (
			imported, skipped, excluded, renamed int
		)

		keyPath := getKey("import")
		storagePath := getStorage("import")

		// Auto-detect format from file extension if --format not set
		if importFormat == "" {
			switch {
			case strings.HasSuffix(importFilePath, ".env"):
				importFormat = "env"
			case strings.HasSuffix(importFilePath, ".json"):
				importFormat = "json"
			case strings.HasSuffix(importFilePath, ".yaml"), strings.HasSuffix(importFilePath, ".yml"):
				importFormat = "yaml"
			default:
				return fmt.Errorf("cannot auto-detect input format from file extension. Please use --format")
			}
		}

		content, err := os.ReadFile(importFilePath)
		if err != nil {
			return fmt.Errorf("failed to read file: %w", err)
		}

		var data map[string]string

		importFormat = strings.ToLower(importFormat)
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

		// Exclude keys if --exclude is set
		if len(importExclude) > 0 {
			excludeMap := make(map[string]bool)
			for _, name := range importExclude {
				excludeMap[name] = true
			}
			for name := range data {
				if excludeMap[name] {
					delete(data, name)
					excluded++
					fmt.Printf("⛔ Excluding key '%s' as per --exclude\n", name)
				}
			}
		}

		key, err := crypto.LoadKeyFromFile(keyPath)
		if err != nil {
			return fmt.Errorf("failed to load key: %w", err)
		}

		store := storage.NewFileStore(storagePath)
		imported = 0

		importRenameMap = make(map[string]string)
		for _, entry := range importRenameList {
			parts := strings.SplitN(entry, "=", 2)
			if len(parts) != 2 {
				return fmt.Errorf("invalid --rename value: %s (expected old=new)", entry)
			}
			importRenameMap[parts[0]] = parts[1]
		}

		for origName, value := range data {
			// Rename key if requested
			if newName, ok := importRenameMap[origName]; ok {
				origName = newName
				renamed++
			}
			name := importPrefix + origName

			// Check existence
			if importSkipExisting && !importForce {
				existing, _ := store.Get(name)
				if existing != nil {
					fmt.Printf("⚠️  Skipping existing secret: %s\n", name)
					skipped++
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

		summary := struct {
			File     string `json:"file"`
			DryRun   bool   `json:"dry_run"`
			Imported int    `json:"imported"`
			Skipped  int    `json:"skipped"`
			Excluded int    `json:"excluded"`
			Renamed  int    `json:"renamed"`
		}{
			File:     importFilePath,
			DryRun:   importDryRun,
			Imported: imported,
			Skipped:  skipped,
			Excluded: excluded,
			Renamed:  renamed,
		}

		if importQuiet {
			return nil
		}

		switch importSummaryFormat {
		case "json":
			out, err := json.MarshalIndent(summary, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to encode summary: %w", err)
			}
			fmt.Println(string(out))
		default:
			if importDryRun {
				log.Printf("✅ Dry run complete. %d secrets would be imported from %s", imported, importFilePath)
			} else {
				log.Printf("✅ Import complete from %s", importFilePath)
			}

			fmt.Println("\n📊 Import Summary")
			fmt.Println("-----------------")
			fmt.Printf("🔐 Imported     : %d\n", summary.Imported)
			fmt.Printf("⚠️  Skipped      : %d\n", summary.Skipped)
			fmt.Printf("⛔ Excluded     : %d\n", summary.Excluded)
			fmt.Printf("✏️  Renamed      : %d\n", summary.Renamed)
			fmt.Println()
		}

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
	importCmd.Flags().StringVar(&importPrefix, "prefix", "", "Optional prefix to apply to all secret names")
	importCmd.Flags().StringSliceVar(&importRenameList, "rename", nil, "Rename secret keys: use --rename old=new")
	importCmd.Flags().StringSliceVar(&importExclude, "exclude", nil, "List of keys to exclude from import")
	importCmd.Flags().StringVar(&importSummaryFormat, "summary", "text", "Output summary format: text (default), json")
	importCmd.Flags().BoolVar(&importQuiet, "quiet", false, "Suppress summary output")

	if err := viper.BindPFlag("import.key", importCmd.Flags().Lookup("key")); err != nil {
		log.Fatalf("Failed to bind config key: %v", err)
	}
	if err := viper.BindPFlag("import.storage", importCmd.Flags().Lookup("storage")); err != nil {
		log.Fatalf("Failed to bind config key: %v", err)
	}

	if err := importCmd.MarkFlagRequired("file"); err != nil {
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
