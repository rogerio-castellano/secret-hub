package cmd

import (
	"encoding/base64"
	"fmt"
	"log"

	"github.com/rogerio-castellano/secret-hub/internal/crypto"
	"github.com/rogerio-castellano/secret-hub/internal/iox"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	inputPath    string
	outputPath   string
	base64Output bool
)

var encryptCmd = &cobra.Command{
	Use:   "encrypt",
	Short: "Encrypt a secret using AES-256-GCM",
	Long: `Encrypt a secret using AES-256-GCM.

This command encrypts a secret from a file or standard input using the AES-256-GCM algorithm. 
You must specify both an input file and an output destination. By default, the encryption key is loaded from 
a configuration file named .secret-hub.yaml in the $HOME directory. If this file is not present, 
a separate key file containing a valid 32-byte encryption key must be provided.

Optionally, the encrypted output can be base64-encoded before being written to the specified output path.

Usage:
  encrypt --input <plaintext> --output <ciphertext> [--key <keyfile>] [--base64]

Examples:
  secret-hub encrypt --in secret.txt --out secret.enc --key mykey.bin
  secret-hub encrypt --in secret.txt --out secret.enc
  secret-hub encrypt --in secret.txt --out secret.enc --key mykey.bin --base64
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		keyPath := getKey("encrypt")
		fmt.Println("Using", keyPath, "to encrypt.")
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

	encryptCmd.Flags().StringVarP(&inputPath, "in", "i", "", "Input file to encrypt (required)")
	encryptCmd.Flags().StringVarP(&outputPath, "out", "o", "", "Output file for encrypted data (required)")
	encryptCmd.Flags().StringP("key", "k", "", "Encryption key path (required unless specified in config).")
	if err := viper.BindPFlag("encrypt.key", encryptCmd.Flags().Lookup("key")); err != nil {
		log.Fatalf("Failed to mark flag required: %v", err)
	}
	encryptCmd.Flags().BoolVar(&base64Output, "base64", false, "Output as base64 instead of raw bytes")

	if err := encryptCmd.MarkFlagRequired("in"); err != nil {
		log.Fatalf("Failed to mark flag required: %v", err)
	}
	if err := encryptCmd.MarkFlagRequired("out"); err != nil {
		log.Fatalf("Failed to mark flag required: %v", err)
	}
}
