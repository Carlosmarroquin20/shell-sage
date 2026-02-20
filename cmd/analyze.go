package cmd

import (
	"fmt"
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/shell-sage/internal/ollama"
	"github.com/spf13/cobra"
)

var analyzeCmd = &cobra.Command{
	Use:   "analyze [file]",
	Short: "Analyze an error log file",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		filePath := args[0]
		content, err := os.ReadFile(filePath)
		if err != nil {
			fmt.Printf("Error reading file: %v\n", err)
			return
		}

		fmt.Printf("🧐 Analyzing %s...\n", filePath)

		// Truncate if too long to avoid token limits (rudimentary handling)
		logContent := string(content)
		if len(logContent) > 2000 {
			logContent = logContent[:2000] + "\n...[truncated]..."
		}

		client := ollama.NewClient()
		prompt := fmt.Sprintf("You are a system administrator. Analyze this log file snippet and summarize the critical errors: \n\n%s", logContent)

		response, err := client.Generate(prompt)
		if err != nil {
			fmt.Printf("Error communicating with Ollama: %v\n", err)
			return
		}

		// Style output
		var style = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#98C379")). // Green-ish for analysis
			PaddingTop(1).
			PaddingBottom(1).
			PaddingLeft(2).
			PaddingRight(2).
			Width(80).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("62"))

		fmt.Println(style.Render(response))
	},
}

func init() {
	rootCmd.AddCommand(analyzeCmd)
}
