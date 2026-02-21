package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/atotto/clipboard"
	"github.com/charmbracelet/lipgloss"
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
			logger.Log.WithError(err).Error("Failed to read shell history")
			metrics.Record("fix", elapsed, err.Error())
			fmt.Println(ui.ErrorStyle().Render("❌ " + err.Error()))
			return
		}

		if len(commands) == 0 {
			fmt.Println(ui.ErrorStyle().Render("⚠️  No recent commands found in history."))
			return
		}

		logger.Log.WithField("commands_found", len(commands)).Info("Shell history read")

		lang := "English"
		if LangFlag != "" {
			lang = LangFlag
		}
		prompt := fmt.Sprintf(
			"IMPORTANT: You MUST respond ONLY in %s. Do not use any other language.\nYou are a shell expert. Given these recent commands, identify if the last one likely failed and suggest a concise fix in max 3 bullet points. Commands: %s",
			lang, strings.Join(commands, " | "),
		)

		client := ollama.NewClient(ModelFlag)

		sp := spinner.New("Scanning history for errors...")
		sp.Start()
		firstToken := true

		borderColor := lipgloss.Color(ui.ColorOrange)
		header := ui.HeaderStyle(ui.ColorOrange).Render("🔧 FIX SUGGESTION")

		response, err := client.GenerateStream(prompt, func(token string) {
			if firstToken {
				sp.Stop()
				firstToken = false
				fmt.Println(header)
				fmt.Println(lipgloss.NewStyle().Foreground(borderColor).Render("╔" + strings.Repeat("═", 76) + "╗"))
				fmt.Print(lipgloss.NewStyle().Foreground(borderColor).Render("║") + "  ")
			}
			formatted := strings.ReplaceAll(token, "\n", "\n"+lipgloss.NewStyle().Foreground(borderColor).Render("║")+"  ")
			fmt.Print(formatted)
		})

		if firstToken {
			sp.Stop()
		}

		elapsed := time.Since(start)

		if err != nil {
			if !firstToken {
				fmt.Println()
			}
			logger.Log.WithError(err).Error("'fix' command failed")
			metrics.Record("fix", elapsed, err.Error())
			fmt.Println(ui.ErrorStyle().Render("❌ " + err.Error()))
			return
		}

		fmt.Println()
		fmt.Println(lipgloss.NewStyle().Foreground(borderColor).Render("╚" + strings.Repeat("═", 76) + "╝"))

		logger.Log.WithField("duration_ms", elapsed.Milliseconds()).Info("'fix' command completed")
		metrics.Record("fix", elapsed, "")

		// Offer clipboard copy
		fmt.Print("\n📋 Copy suggestion to clipboard? [y/N]: ")
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(strings.ToLower(input))
		if input == "y" || input == "yes" {
			if err := clipboard.WriteAll(response); err != nil {
				logger.Log.WithError(err).Warn("Failed to copy to clipboard")
				fmt.Println(ui.ErrorStyle().Render("❌ Could not copy: " + err.Error()))
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
