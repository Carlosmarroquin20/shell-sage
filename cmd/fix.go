package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/charmbracelet/lipgloss"
	"github.com/shell-sage/internal/history"
	"github.com/shell-sage/internal/ollama"
	"github.com/shell-sage/internal/spinner"
	"github.com/spf13/cobra"
)

var fixCmd = &cobra.Command{
	Use:   "fix",
	Short: "Analyze recent shell history and suggest a fix for the last error",
	Run: func(cmd *cobra.Command, args []string) {
		// Read last 10 commands from history
		commands, err := history.GetRecentCommands(10)
		if err != nil {
			fmt.Printf("❌ Could not read shell history: %v\n", err)
			return
		}

		if len(commands) == 0 {
			fmt.Println("⚠️  No recent commands found in history.")
			return
		}

		// Start spinner while AI thinks
		sp := spinner.New("Scanning history for errors...")
		sp.Start()

		client := ollama.NewClient(ModelFlag)
		prompt := fmt.Sprintf(
			"You are a shell expert. Given these recent commands, identify if the last one likely failed and suggest a concise fix in max 3 bullet points. Commands: %s",
			strings.Join(commands, " | "),
		)

		response, err := client.Generate(prompt)
		sp.Stop()

		if err != nil {
			fmt.Printf("❌ Error: %v\n", err)
			return
		}

		// Header label
		headerStyle := lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FF8C00")).
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

		// Offer to copy the suggestion to clipboard
		fmt.Print("\n📋 Copy suggestion to clipboard? [y/N]: ")
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(strings.ToLower(input))
		if input == "y" || input == "yes" {
			if err := clipboard.WriteAll(response); err != nil {
				fmt.Println("❌ Could not copy to clipboard:", err)
			} else {
				fmt.Println("✅ Copied to clipboard!")
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(fixCmd)
}
