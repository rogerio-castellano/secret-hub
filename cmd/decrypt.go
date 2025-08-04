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
	decInputPath  string
	decOutputPath string
	base64Input   bool
)

var decryptCmd = &cobra.Command{
	Use:   "decrypt",
	Short: "Decrypt a secret using AES-256-GCM",
	Long: `Decrypt a secret encrypted with AES-256-GCM.

This command decrypts a secret from an encrypted file or standard input using the AES-256-GCM algorithm. 
You must specify both the encrypted input source and the destination for the decrypted output. By default, 
the decryption key is loaded from a .secret-hub.yaml configuration file in the $HOME directory. If this 
file is unavailable, you must provide a separate key file containing a valid 32-byte encryption key.

If the encrypted input is base64-encoded, enable decoding before decryption. The result will be written to 
the specified output file.

Usage:
  decrypt --input <ciphertext> --output <plaintext> [--key <keyfile>] [--base64]
  
Examples:
# Decrypt 'secret.enc' using the key from 'key.bin'; output saved as 'secret-dec.txt'
secret-hub decrypt --in secret.enc --out secret-dec.txt --key key.bin

# Decrypt 'secret.enc' using default key and settings from the configuration file
secret-hub decrypt --in secret.enc --out secret-dec.txt

# Decrypt Base64-encoded 'secret.enc' using 'key.bin'; decoded output saved as 'secret-dec.txt'
secret-hub decrypt --in secret.enc --out secret-dec.txt --key key.bin --base64
  `,
	RunE: func(cmd *cobra.Command, args []string) error {
		keyPath := getKey("decrypt")
		key, err := crypto.LoadKeyFromFile(keyPath)
		log.Println("🔑 Loading decryption key...", keyPath)
		if err != nil {
			return fmt.Errorf("failed to load key: %w", err)
		}

		ciphertext, err := iox.ReadInput(decInputPath)
		if err != nil {
			return fmt.Errorf("failed to read input: %w", err)
		}

		if base64Input {
			ciphertext, err = base64.StdEncoding.DecodeString(string(ciphertext))
			if err != nil {
				return fmt.Errorf("base64 decoding failed(%s): %w", string(ciphertext), err)
			}
		}

		plaintext, err := crypto.Decrypt(key, ciphertext)
		if err != nil {
			return fmt.Errorf("decryption failed: %w", err)
		}

		if err := iox.WriteOutput(decOutputPath, plaintext); err != nil {
			return fmt.Errorf("failed to write output: %w", err)
		}

		log.Printf("🔓 Secret decrypted successfully.")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(decryptCmd)

	decryptCmd.Flags().StringVarP(&decInputPath, "in", "i", "", "Encrypted input file (required)")
	decryptCmd.Flags().StringVarP(&decOutputPath, "out", "o", "", "Decrypted output file (required)")
	decryptCmd.Flags().StringP("key", "k", "", "Decryption key path (required unless specified in config).")

	if err := viper.BindPFlag("decrypt.key", decryptCmd.Flags().Lookup("key")); err != nil {
		log.Fatalf("Failed to bind config key: %v", err)
	}

	decryptCmd.Flags().BoolVar(&base64Input, "base64", false, "Input is base64 encoded")

	if err := decryptCmd.MarkFlagRequired("in"); err != nil {
		log.Fatalf("Failed to mark flag required: %v", err)
	}
	if err := decryptCmd.MarkFlagRequired("out"); err != nil {
		log.Fatalf("Failed to mark flag required: %v", err)
	}
}
