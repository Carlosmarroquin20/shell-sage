package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/atotto/clipboard"
	"github.com/shell-sage/internal/history"
	"github.com/shell-sage/internal/logger"
	"github.com/shell-sage/internal/metrics"
	"github.com/shell-sage/internal/ollama"
	"github.com/shell-sage/internal/spinner"
	"github.com/shell-sage/internal/ui"
	"github.com/spf13/cobra"
)

var fixCmd = &cobra.Command{
	Use:   "fix",
	Short: "Analyze recent shell history and suggest a fix for the last error",
	Run: func(cmd *cobra.Command, args []string) {
		start := time.Now()

		logger.Log.Info("Starting 'fix' command")

		commands, err := history.GetRecentCommands(10)
		if err != nil {
			elapsed := time.Since(start)
			logger.Log.WithError(err).WithField("duration_ms", elapsed.Milliseconds()).Error("Failed to read shell history")
			metrics.Record("fix", elapsed, err.Error())
			fmt.Println(ui.ErrorStyle().Render("❌ " + err.Error()))
			return
		}

		logger.Log.WithField("commands_found", len(commands)).Info("Shell history read successfully")

		if len(commands) == 0 {
			fmt.Println(ui.ErrorStyle().Render("⚠️  No recent commands found in history."))
			return
		}

		sp := spinner.New("Scanning history for errors...")
		sp.Start()

		client := ollama.NewClient(ModelFlag)
		prompt := fmt.Sprintf(
			"You are a shell expert. Given these recent commands, identify if the last one likely failed and suggest a concise fix in max 3 bullet points. Commands: %s",
			strings.Join(commands, " | "),
		)

		response, err := client.Generate(prompt)
		sp.Stop()

		elapsed := time.Since(start)

		if err != nil {
			logger.Log.WithError(err).WithField("duration_ms", elapsed.Milliseconds()).Error("'fix' command failed during AI generation")
			metrics.Record("fix", elapsed, err.Error())
			fmt.Println(ui.ErrorStyle().Render("❌ " + err.Error()))
			return
		}

		logger.Log.WithField("duration_ms", elapsed.Milliseconds()).Info("'fix' command completed successfully")
		metrics.Record("fix", elapsed, "")

		header := ui.HeaderStyle(ui.ColorOrange).Render("🔧 FIX SUGGESTION")
		fmt.Println(header)
		fmt.Println(ui.BodyStyleDouble(ui.ColorOrange).Render(response))

		// Offer to copy to clipboard
		fmt.Print("\n📋 Copy suggestion to clipboard? [y/N]: ")
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(strings.ToLower(input))
		if input == "y" || input == "yes" {
			if err := clipboard.WriteAll(response); err != nil {
				logger.Log.WithError(err).Warn("Failed to copy to clipboard")
				fmt.Println(ui.ErrorStyle().Render("❌ Could not copy to clipboard: " + err.Error()))
			} else {
				logger.Log.Info("Response copied to clipboard")
				fmt.Println("✅ Copied to clipboard!")
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(fixCmd)
}
