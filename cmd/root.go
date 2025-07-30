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
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
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
