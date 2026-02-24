package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/atotto/clipboard"
	"github.com/charmbracelet/lipgloss"
	"github.com/shell-sage/internal/logger"
	"github.com/shell-sage/internal/metrics"
	"github.com/shell-sage/internal/ollama"
	"github.com/shell-sage/internal/spinner"
	"github.com/shell-sage/internal/ui"
	"github.com/spf13/cobra"
)

const maxLogChars = 2000

var analyzeCmd = &cobra.Command{
	Use:   "analyze [file]",
	Short: "Analyze an error log file and summarize critical issues",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		start := time.Now()
		filePath := args[0]

		logger.Log.WithField("file", filePath).Info("Starting 'analyze' command")

		content, err := os.ReadFile(filePath)
		if err != nil {
			elapsed := time.Since(start)
			logger.Log.WithError(err).Error("Failed to read log file")
			metrics.Record("analyze", elapsed, err.Error())
			fmt.Println(ui.ErrorStyle().Render("❌ Error reading file: " + err.Error()))
			return
		}

		logContent := string(content)
		fullSize := len(logContent)
		logger.Log.WithField("file_size_chars", fullSize).Info("Log file read")

		if fullSize > maxLogChars {
			fmt.Printf("⚠️  Log file is large (%d chars). Send full content to AI? This may be slow. [y/N]: ", fullSize)
			reader := bufio.NewReader(os.Stdin)
			input, _ := reader.ReadString('\n')
			if strings.TrimSpace(strings.ToLower(input)) != "y" {
				logContent = logContent[:maxLogChars] + "\n...[truncated]..."
				logger.Log.WithField("truncated_at", maxLogChars).Info("Log truncated by user choice")
				fmt.Println("📄 Using first 2000 characters.")
			} else {
				logger.Log.Info("User chose to send full log")
				fmt.Println("📄 Sending full log to AI...")
			}
		}

		lang := "English"
		if LangFlag != "" {
			lang = LangFlag
		}
		prompt := fmt.Sprintf(
			"IMPORTANT: You MUST respond ONLY in %s. Do not use any other language.\nYou are a sysadmin. Analyze this log and summarize the critical errors in max 4 bullet points, no intro:\n\n%s",
			lang, logContent,
		)

		client := ollama.NewClient(ModelFlag)

		sp := spinner.New(fmt.Sprintf("Analyzing %s...", filePath))
		sp.Start()
		firstToken := true

		borderColor := lipgloss.Color(ui.ColorGreen)
		header := ui.HeaderStyle(ui.ColorGreen).Render("🧠 LOG ANALYSIS › " + filePath)

		response, err := client.GenerateStream(prompt, func(token string) {
			if firstToken {
				sp.Stop()
				firstToken = false
				fmt.Println(header)
				fmt.Println(lipgloss.NewStyle().Foreground(borderColor).Render("╭" + strings.Repeat("─", 76) + "╮"))
				fmt.Print(lipgloss.NewStyle().Foreground(borderColor).Render("│") + "  ")
			}
			formatted := strings.ReplaceAll(token, "\n", "\n"+lipgloss.NewStyle().Foreground(borderColor).Render("│")+"  ")
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
			logger.Log.WithError(err).Error("'analyze' command failed")
			metrics.Record("analyze", elapsed, err.Error())
			fmt.Println(ui.ErrorStyle().Render("❌ " + err.Error()))
			return
		}

		fmt.Println()
		fmt.Println(lipgloss.NewStyle().Foreground(borderColor).Render("╰" + strings.Repeat("─", 76) + "╯"))

		logger.Log.WithField("duration_ms", elapsed.Milliseconds()).Info("'analyze' command completed")
		metrics.Record("analyze", elapsed, "")

		if CopyFlag {
			if err := clipboard.WriteAll(response); err != nil {
				logger.Log.WithError(err).Warn("Failed to copy to clipboard")
				fmt.Println(ui.ErrorStyle().Render("\n❌ Could not copy: " + err.Error()))
			} else {
				fmt.Println("\n✅ Copied to clipboard!")
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(analyzeCmd)
}
