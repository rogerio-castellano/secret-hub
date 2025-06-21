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
	decInputPath  string
	decOutputPath string
	decKeyPath    string
	base64Input   bool
)

var decryptCmd = &cobra.Command{
	Use:   "decrypt",
	Short: "Decrypt a secret using AES-256-GCM",
	Long: `The "decrypt" command allows you to decrypt a secret that was previously encrypted using the AES-256-GCM algorithm. 
You must provide a valid decryption key file and the encrypted input (either as a file or from standard input). 
Optionally, if the input is base64-encoded, you can specify this to decode before decryption. 
The decrypted plaintext will be written to the specified output path or to standard output if no path is provided.

Usage:
  decrypt --key <keyfile> --input <ciphertext> [--output <plaintext>] [--base64]`,
	RunE: func(cmd *cobra.Command, args []string) error {
		key, err := crypto.LoadKeyFromFile(decKeyPath)
		log.Println("🔑 Loading decryption key...", decKeyPath)
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

	decryptCmd.Flags().StringVarP(&decInputPath, "in", "i", "", "Encrypted input file")
	decryptCmd.Flags().StringVarP(&decOutputPath, "out", "o", "", "Decrypted output file")
	decryptCmd.Flags().StringVarP(&decKeyPath, "key", "k", "", "Path to 32-byte decryption key")
	decryptCmd.Flags().BoolVar(&base64Input, "base64", false, "Input is base64 encoded")

	decryptCmd.MarkFlagRequired("in")
	decryptCmd.MarkFlagRequired("out")
	decryptCmd.MarkFlagRequired("key")
}
