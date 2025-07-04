package cmd

import (
	"encoding/base64"
	"fmt"
	"log"

	"github.com/rogerio-castellano/secret-hub/internal/crypto"
	"github.com/rogerio-castellano/secret-hub/internal/iox"
	"github.com/spf13/cobra"
)

var (
	inputPath    string
	outputPath   string
	keyPath      string
	base64Output bool
)

var encryptCmd = &cobra.Command{
	Use:   "encrypt",
	Short: "Encrypt a secret using AES-256-GCM",
	Long: `Encrypts a secret file using AES-256-GCM encryption.

You must provide an input file, an output file, and a 32-byte key file.
Optionally, you can output the encrypted data as base64.

Example:
  secret-hub encrypt --in secret.txt --out secret.enc --key mykey.bin --base64
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Load key
		key, err := crypto.LoadKeyFromFile(keyPath)
		if err != nil {
			return fmt.Errorf("failed to load key: %w", err)
		}

		// Read input
		plaintext, err := iox.ReadInput(inputPath)
		if err != nil {
			return fmt.Errorf("failed to read input: %w", err)
		}

		// Encrypt
		ciphertext, err := crypto.Encrypt(key, plaintext)
		if err != nil {
			return fmt.Errorf("encryption failed: %w", err)
		}

		// Write output
		if base64Output {
			encoded := base64.StdEncoding.EncodeToString(ciphertext)
			if err := iox.WriteOutput(outputPath, []byte(encoded)); err != nil {
				return fmt.Errorf("failed to write base64 output: %w", err)
			}
		} else {
			if err := iox.WriteOutput(outputPath, ciphertext); err != nil {
				return fmt.Errorf("failed to write output file: %w", err)
			}

		}

		log.Println("🔒 Secret encrypted successfully.", outputPath)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(encryptCmd)

	encryptCmd.Flags().StringVarP(&inputPath, "in", "i", "", "Input file to encrypt")
	encryptCmd.Flags().StringVarP(&outputPath, "out", "o", "", "Output file for encrypted data")
	encryptCmd.Flags().StringVarP(&keyPath, "key", "k", "", "Path to 32-byte encryption key")
	encryptCmd.Flags().BoolVar(&base64Output, "base64", false, "Output as base64 instead of raw bytes")

	encryptCmd.MarkFlagRequired("input")
	encryptCmd.MarkFlagRequired("output")
	encryptCmd.MarkFlagRequired("key")
}
