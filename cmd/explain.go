package cmd

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/shell-sage/internal/ollama"
	"github.com/spf13/cobra"
)

var explainCmd = &cobra.Command{
	Use:   "explain [command]",
	Short: "Explain a shell command",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		commandToExplain := strings.Join(args, " ")
		
		fmt.Printf("🤔 Asking AI to explain: %s...\n", commandToExplain)

		client := ollama.NewClient()
		prompt := fmt.Sprintf("You are an expert in Linux/macOS terminals. Explain briefly what this command does using bullet points for flags: '%s'. Keep it concise.", commandToExplain)

		response, err := client.Generate(prompt)
		if err != nil {
			fmt.Printf("Error communicating with Ollama: %v\n", err)
			return
		}

		// Styling output with lipgloss
		var style = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			PaddingTop(1).
			PaddingBottom(1).
			PaddingLeft(2).
			PaddingRight(2).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("238"))

		fmt.Println(style.Render(response))
	},
}

func init() {
	rootCmd.AddCommand(explainCmd)
}
