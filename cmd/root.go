package cmd

import (
	"os"
	"time"

	"github.com/shell-sage/internal/config"
	"github.com/shell-sage/internal/pipeline"
	"github.com/shell-sage/internal/pipeline/middleware/cache"
	"github.com/shell-sage/internal/pipeline/middleware/enhancer"
	"github.com/shell-sage/internal/pipeline/middleware/retry"
	"github.com/shell-sage/internal/provider"

	"github.com/spf13/cobra"
)

// ModelFlag holds the value of the --model flag, available to all subcommands.
var ModelFlag string

// LangFlag holds the value of the --lang flag (e.g. "es", "fr", "English").
var LangFlag string

// CopyFlag determines if the output should be copied to the clipboard.
var CopyFlag bool

// ProviderFlag holds the value of the --provider flag (e.g. "ollama").
// Priority at runtime: flag > SSAGE_PROVIDER env > config file > "ollama".
var ProviderFlag string

var rootCmd = &cobra.Command{
	Use:   "ssage",
	Short: "Shell Sage - Your AI Terminal Assistant",
	Long: `Shell Sage is a CLI tool that uses local AI to help you
explain commands, fix errors from your history, and analyze logs.

Providers are pluggable: use --provider to select a backend (default: ollama).`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Load config and use as fallback for flags not explicitly set.
		cfg, err := config.Load()
		if err != nil {
			return // Silently ignore config errors for regular commands
		}

		if LangFlag == "" && cfg.Lang != "" {
			LangFlag = cfg.Lang
		}
		if ProviderFlag == "" && cfg.Provider != "" {
			ProviderFlag = cfg.Provider
		}
		// Model priority (flag > env > config) is handled inside provider.New
		// → ollama.NewClient, keeping the logic co-located with the backend.
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	// Persistent flags available to all subcommands.
	rootCmd.PersistentFlags().StringVarP(&ModelFlag, "model", "m", "", "LLM model to use (e.g. llama3, mistral)")
	rootCmd.PersistentFlags().StringVarP(&LangFlag, "lang", "l", "", "Response language (e.g. 'es' for Spanish, 'fr' for French)")
	rootCmd.PersistentFlags().BoolVarP(&CopyFlag, "copy", "c", false, "Copy the suggested command or explanation to clipboard")
	rootCmd.PersistentFlags().StringVarP(&ProviderFlag, "provider", "p", "", "AI provider to use (e.g. ollama)")
}

// buildPipeline creates a ready-to-use Pipeline wired with the standard
// middleware stack: enhancer → cache → retry → provider.
//
// The middleware order ensures that:
//  1. enhancer runs first to inject OS/Shell context into the prompt.
//  2. cache uses the enhanced prompt as its key and short-circuits on a hit.
//  3. retry wraps the actual provider call to handle transient errors.
func buildPipeline() (*pipeline.Pipeline, error) {
	p, err := provider.New(ProviderFlag, ModelFlag)
	if err != nil {
		return nil, err
	}
	return pipeline.New(
		p,
		enhancer.New(),
		cache.New(24*time.Hour, "tip"),
		retry.New(3),
	), nil
}
