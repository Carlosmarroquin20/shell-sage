package cmd

import (
	"os"

	"github.com/shell-sage/internal/config"

	"github.com/spf13/cobra"
)

// ModelFlag holds the value of the --model flag, available to all subcommands.
var ModelFlag string

// LangFlag holds the value of the --lang flag (e.g. "es", "fr", "English").
var LangFlag string

// CopyFlag determines if the output should be copied to the clipboard.
var CopyFlag bool

var rootCmd = &cobra.Command{
	Use:   "ssage",
	Short: "Shell Sage - Your AI Terminal Assistant",
	Long: `Shell Sage is a CLI tool that uses local AI (Ollama) to help you 
explain commands, fix errors from your history, and analyze logs.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Load config and use as fallback for flags
		cfg, err := config.Load()
		if err != nil {
			return // Silently ignore config errors for regular commands
		}

		if LangFlag == "" && cfg.Lang != "" {
			LangFlag = cfg.Lang
		}
		// Model is handled inside NewClient for priority flag > env > config
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	// Persistent flags available to all subcommands
	rootCmd.PersistentFlags().StringVarP(&ModelFlag, "model", "m", "", "Ollama model to use (e.g. llama3, mistral)")
	rootCmd.PersistentFlags().StringVarP(&LangFlag, "lang", "l", "", "Response language (e.g. 'es' for Spanish, 'fr' for French)")
	rootCmd.PersistentFlags().BoolVarP(&CopyFlag, "copy", "c", false, "Copy the suggested command or explanation to clipboard")
}
