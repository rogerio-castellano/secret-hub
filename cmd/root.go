/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "secret-hub",
	Short: "CLI tool for encrypting, storing, and retrieving secrets securely",
	Long: `secret-hub is a command-line tool written in Go for managing secrets locally using AES-256 encryption.

It supports:
- Generating secure keys
- Encrypting and decrypting files
- Storing encrypted secrets by name
- Retrieving, listing, searching, and deleting secrets

Secrets can be supplied provided either as direct key values or sourced from external files. 
The resulting data is saved locally for easy retrieval. 
Keys are structured in JSON format, and all files are securely stored using encryption to protect sensitive information.`,
}

// Execute runs the root command
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.secret-hub.yaml)")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Enable verbose output")
	if err := viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose")); err != nil {
		log.Fatalf("Failed to bind config key: %v", err)
	}

	viper.SetDefault("storage", "secrets.json")
}

func initConfig() {
	// If a config file is found, read it in.
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)
		viper.AddConfigPath(home)
		viper.SetConfigName(".secret-hub")
		viper.SetConfigType("yaml")
	}

	viper.AutomaticEnv() // load env vars

	if err := viper.ReadInConfig(); err == nil {
		log.Println("📘 Using config file:", viper.ConfigFileUsed())
	} else {
		log.Println("⚠️ No config file found (optional):", err)
	}
}

func getKey(namespace string) string {
	val := viper.GetString(namespace + ".key")
	if val == "" {
		val = viper.GetString("key")
	}
	return val
}

func getStorage(namespace string) string {
	val := viper.GetString(namespace + ".storage")
	if val == "" {
		val = viper.GetString("storage")
	}
	return val
}
