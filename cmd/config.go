package cmd

import (
	"fmt"

	"github.com/shell-sage/internal/config"

	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage Shell Sage configuration",
	Long:  `View or set persistent configuration like default model and language.`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Load()
		if err != nil {
			fmt.Printf("Error loading config: %v\n", err)
			return
		}
		path, _ := config.GetConfigPath()
		fmt.Printf("Config file: %s\n", path)
		fmt.Printf("Model: %s\n", cfg.Model)
		fmt.Printf("Language: %s\n", cfg.Lang)
		fmt.Printf("Provider: %s\n", cfg.Provider)
	},
}

var setConfigCmd = &cobra.Command{
	Use:   "set [key] [value]",
	Short: "Set a configuration value",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Load()
		if err != nil {
			fmt.Printf("Error loading config: %v\n", err)
			return
		}

		key := args[0]
		value := args[1]

		switch key {
		case "model":
			cfg.Model = value
		case "lang":
			cfg.Lang = value
		case "provider":
			cfg.Provider = value
		default:
			fmt.Printf("Unknown config key: %s (available: model, lang, provider)\n", key)
			return
		}

		if err := cfg.Save(); err != nil {
			fmt.Printf("Error saving config: %v\n", err)
			return
		}

		fmt.Printf("Successfully set %s to %s\n", key, value)
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(setConfigCmd)
}
