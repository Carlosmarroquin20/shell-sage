package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// ModelFlag holds the value of the --model flag, available to all subcommands.
var ModelFlag string

var rootCmd = &cobra.Command{
	Use:   "ssage",
	Short: "Shell Sage - Your AI Terminal Assistant",
	Long: `Shell Sage is a CLI tool that uses local AI (Ollama) to help you 
explain commands, fix errors from your history, and analyze logs.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	// Persistent flag available to all subcommands
	rootCmd.PersistentFlags().StringVarP(&ModelFlag, "model", "m", "", "Ollama model to use (e.g. llama3, mistral)")
}
