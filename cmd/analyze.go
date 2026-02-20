package cmd

import (
	"fmt"
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/shell-sage/internal/ollama"
	"github.com/shell-sage/internal/spinner"
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
			fmt.Printf("❌ Error reading file: %v\n", err)
			return
		}

		// Start spinner while AI thinks
		sp := spinner.New(fmt.Sprintf("Analyzing %s...", filePath))
		sp.Start()

		// Truncate to avoid token limits
		logContent := string(content)
		if len(logContent) > 2000 {
			logContent = logContent[:2000] + "\n...[truncated]..."
		}

		client := ollama.NewClient(ModelFlag)
		prompt := fmt.Sprintf("You are a sysadmin. Analyze this log and summarize the critical errors in max 4 bullet points, no intro: \n\n%s", logContent)

		response, err := client.Generate(prompt)
		sp.Stop()

		if err != nil {
			fmt.Printf("❌ Error: %v\n", err)
			return
		}

		// Header label
		headerStyle := lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#39FF14")).
			Background(lipgloss.Color("#1a1a2e")).
			Padding(0, 1)

		header := headerStyle.Render("🧠 LOG ANALYSIS › " + filePath)

		// Body style
		bodyStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#E0E0E0")).
			PaddingTop(1).
			PaddingBottom(1).
			PaddingLeft(2).
			PaddingRight(2).
			Width(76).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#39FF14"))

		fmt.Println(header)
		fmt.Println(bodyStyle.Render(response))
	},
}

func init() {
	rootCmd.AddCommand(analyzeCmd)
}
