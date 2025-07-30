// cmd/genkey.go
package cmd

import (
	"crypto/rand"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	keyOutputPath string
)

var generateKeyCmd = &cobra.Command{
	Use:   "generate-key",
	Short: "Generate a 256-bit AES key",
	Long: `Generate a secure 256-bit AES key for encryption and decryption.

This command creates a cryptographically strong 32-byte (256-bit) random key suitable for AES-256-GCM operations. The key is generated using a secure random source and saved to a specified file path. This key can later be used with other commands such as encrypt, decrypt, or store.

By default, the key is written to a file named key.bin unless a different path is provided using the --out flag. The output file is saved with restrictive permissions to safeguard against unauthorized access.

Usage 
	generate-key [--out <keyfile>]

Example:
  secret-hub generate-key --out custom-key.bin
  secret-hub generate-key
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		key := make([]byte, 32)
		_, err := rand.Read(key)
		if err != nil {
			return fmt.Errorf("failed to generate key: %w", err)
		}

		err = os.WriteFile(keyOutputPath, key, 0600)
		if err != nil {
			return fmt.Errorf("failed to write key to file: %w", err)
		}

		fmt.Printf("🔑 256-bit key generated and written to %s\n", keyOutputPath)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(generateKeyCmd)

	generateKeyCmd.Flags().StringVarP(&keyOutputPath, "out", "o", "key.bin", "Path to write the generated key")
}
