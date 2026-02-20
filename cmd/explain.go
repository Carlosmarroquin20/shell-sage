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
		prompt := fmt.Sprintf("Explain this shell command in max 3 bullet points. Be extremely concise, no intro, no extra text: '%s'", commandToExplain)

		response, err := client.Generate(prompt)
		if err != nil {
			fmt.Printf("Error communicating with Ollama: %v\n", err)
			return
		}

		// Header label
		headerStyle := lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#00D7FF")). // Bright cyan
			Background(lipgloss.Color("#1a1a2e")).
			Padding(0, 1)

		header := headerStyle.Render("⚡ EXPLAIN")

		// Body style — no heavy background, clean border
		bodyStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#E0E0E0")).
			PaddingTop(1).
			PaddingBottom(1).
			PaddingLeft(2).
			PaddingRight(2).
			Width(76).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#00D7FF"))

		fmt.Println(header)
		fmt.Println(bodyStyle.Render(response))
	},
}

func init() {
	rootCmd.AddCommand(explainCmd)
}
