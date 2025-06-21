/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
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

// decryptCmd represents the decrypt command
var decryptCmd = &cobra.Command{
	Use:   "decrypt",
	Short: "Decrypt a secret using AES-256-GCM",
	//TODO: Add a Long description
	// 	Long: `A longer description that spans multiple lines and likely contains examples
	// and usage of using your command. For example:

	// Cobra is a CLI library for Go that empowers applications.
	// This application is a tool to generate the needed files
	// to quickly create a Cobra application.`,
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

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// decryptCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// decryptCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
