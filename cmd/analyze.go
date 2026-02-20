package cmd

import (
	"fmt"
	"os"

	"github.com/shell-sage/internal/ollama"
	"github.com/shell-sage/internal/spinner"
	"github.com/shell-sage/internal/ui"
	"github.com/spf13/cobra"
)

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

		sp := spinner.New(fmt.Sprintf("Analyzing %s...", filePath))
		sp.Start()

		logContent := string(content)
		if len(logContent) > 2000 {
			logContent = logContent[:2000] + "\n...[truncated]..."
		}

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
