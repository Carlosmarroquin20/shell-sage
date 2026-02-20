package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

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
		filePath := args[0]
		content, err := os.ReadFile(filePath)
		if err != nil {
			fmt.Println(ui.ErrorStyle().Render("❌ Error reading file: " + err.Error()))
			return
		}

		logContent := string(content)

		// If the log is too large, ask the user interactively
		if len(logContent) > maxLogChars {
			fmt.Printf("⚠️  Log file is large (%d chars). Send full content to AI? This may be slow. [y/N]: ", len(logContent))
			reader := bufio.NewReader(os.Stdin)
			input, _ := reader.ReadString('\n')
			if strings.TrimSpace(strings.ToLower(input)) != "y" {
				// Keep only the first maxLogChars characters, focusing on the start (usually where errors start)
				logContent = logContent[:maxLogChars] + "\n...[truncated — run with full content by choosing 'y']..."
				fmt.Println("📄 Using first 2000 characters of the log.")
			} else {
				fmt.Println("📄 Sending full log to AI...")
			}
		}

		sp := spinner.New(fmt.Sprintf("Analyzing %s...", filePath))
		sp.Start()

		client := ollama.NewClient(ModelFlag)
		prompt := fmt.Sprintf("You are a sysadmin. Analyze this log and summarize the critical errors in max 4 bullet points, no intro: \n\n%s", logContent)

		response, err := client.Generate(prompt)
		sp.Stop()

		if err != nil {
			fmt.Println(ui.ErrorStyle().Render("❌ " + err.Error()))
			return
		}

		header := ui.HeaderStyle(ui.ColorGreen).Render("🧠 LOG ANALYSIS › " + filePath)
		fmt.Println(header)
		fmt.Println(ui.BodyStyle(ui.ColorGreen).Render(response))
	},
}

func init() {
	rootCmd.AddCommand(analyzeCmd)
}
