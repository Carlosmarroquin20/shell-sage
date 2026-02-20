package cmd

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/shell-sage/internal/history"
	"github.com/shell-sage/internal/ollama"
	"github.com/spf13/cobra"
)

var fixCmd = &cobra.Command{
	Use:   "fix",
	Short: "Fix the last broken command from history",
	Run: func(cmd *cobra.Command, args []string) {
		// Read last 10 commands
		commands, err := history.GetRecentCommands(10)
		if err != nil {
			fmt.Printf("Error reading history: %v\n", err)
			return
		}

		if len(commands) == 0 {
			fmt.Println("No recent commands found in history.")
			return
		}

		// Provide context to AI
		fmt.Println("🕵️  Scanning history for errors...")

		client := ollama.NewClient()
		prompt := fmt.Sprintf("Analyze these recent shell commands and identify if there is a mistake in the last one, or suggest a fix for a likely failed command. Usage context: %v. Return a concise fix and explanation.", commands)

		response, err := client.Generate(prompt)
		if err != nil {
			fmt.Printf("Error communicating with Ollama: %v\n", err)
			return
		}

		// Header label
		headerStyle := lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FF8C00")). // Orange
			Background(lipgloss.Color("#1a1a2e")).
			Padding(0, 1)

		header := headerStyle.Render("🔧 FIX SUGGESTION")

		// Body style
		bodyStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#E0E0E0")).
			PaddingTop(1).
			PaddingBottom(1).
			PaddingLeft(2).
			PaddingRight(2).
			Width(76).
			BorderStyle(lipgloss.DoubleBorder()).
			BorderForeground(lipgloss.Color("#FF8C00"))

		fmt.Println(header)
		fmt.Println(bodyStyle.Render(response))
	},
}

func init() {
	rootCmd.AddCommand(fixCmd)
}
