package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

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
			logger.Log.WithError(err).WithField("file", filePath).Error("Failed to read log file")
			metrics.Record("analyze", elapsed, err.Error())
			fmt.Println(ui.ErrorStyle().Render("❌ Error reading file: " + err.Error()))
			return
		}

		logContent := string(content)
		fullSize := len(logContent)
		logger.Log.WithField("file_size_chars", fullSize).Info("Log file read")

		// Interactively ask if the file is too large
		if fullSize > maxLogChars {
			fmt.Printf("⚠️  Log file is large (%d chars). Send full content to AI? This may be slow. [y/N]: ", fullSize)
			reader := bufio.NewReader(os.Stdin)
			input, _ := reader.ReadString('\n')
			if strings.TrimSpace(strings.ToLower(input)) != "y" {
				logContent = logContent[:maxLogChars] + "\n...[truncated]..."
				logger.Log.WithField("truncated_at", maxLogChars).Info("Log content truncated by user choice")
				fmt.Println("📄 Using first 2000 characters of the log.")
			} else {
				logger.Log.Info("User chose to send full log content")
				fmt.Println("📄 Sending full log to AI...")
			}
		}

		sp := spinner.New(fmt.Sprintf("Analyzing %s...", filePath))
		sp.Start()

		lang := "English"
		if LangFlag != "" {
			lang = LangFlag
		}
		client := ollama.NewClient(ModelFlag)
		prompt := fmt.Sprintf("IMPORTANT: You MUST respond ONLY in %s. Do not use any other language.\nYou are a sysadmin. Analyze this log and summarize the critical errors in max 4 bullet points, no intro:\n\n%s", lang, logContent)

		response, err := client.Generate(prompt)
		sp.Stop()

		elapsed := time.Since(start)

		if err != nil {
			logger.Log.WithError(err).WithField("duration_ms", elapsed.Milliseconds()).Error("'analyze' command failed during AI generation")
			metrics.Record("analyze", elapsed, err.Error())
			fmt.Println(ui.ErrorStyle().Render("❌ " + err.Error()))
			return
		}

		logger.Log.WithField("duration_ms", elapsed.Milliseconds()).Info("'analyze' command completed successfully")
		metrics.Record("analyze", elapsed, "")

		header := ui.HeaderStyle(ui.ColorGreen).Render("🧠 LOG ANALYSIS › " + filePath)
		fmt.Println(header)
		fmt.Println(ui.BodyStyle(ui.ColorGreen).Render(response))
	},
}

func init() {
	rootCmd.AddCommand(analyzeCmd)
}
