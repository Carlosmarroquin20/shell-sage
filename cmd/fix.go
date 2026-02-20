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

		// Style output
		var style = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#E06C75")). // Red-ish for fix
			PaddingTop(1).
			PaddingBottom(1).
			PaddingLeft(2).
			PaddingRight(2).
			Width(80).
			BorderStyle(lipgloss.DoubleBorder()).
			BorderForeground(lipgloss.Color("208"))

		fmt.Println(style.Render(response))
	},
}

func init() {
	rootCmd.AddCommand(fixCmd)
}
