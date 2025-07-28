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
	Long: `Generate a secure 256-bit (32-byte) random key for use with AES-256 encryption.

Example:
  secret-hub generate-key --out mykey.bin`,
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

	generateKeyCmd.Flags().StringVarP(&keyOutputPath, "out", "o", "mykey.bin", "Path to write the generated key")
}
