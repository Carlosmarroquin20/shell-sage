package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/shell-sage/internal/history"
	"github.com/shell-sage/internal/ollama"
	"github.com/shell-sage/internal/spinner"
	"github.com/shell-sage/internal/ui"
	"github.com/spf13/cobra"
)

var fixCmd = &cobra.Command{
	Use:   "fix",
	Short: "Analyze recent shell history and suggest a fix for the last error",
	Run: func(cmd *cobra.Command, args []string) {
		commands, err := history.GetRecentCommands(10)
		if err != nil {
			fmt.Println(ui.ErrorStyle().Render("❌ " + err.Error()))
			return
		}

		if len(commands) == 0 {
			fmt.Println(ui.ErrorStyle().Render("⚠️  No recent commands found in history."))
			return
		}

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
			fmt.Println(ui.ErrorStyle().Render("❌ " + err.Error()))
			return
		}

		header := ui.HeaderStyle(ui.ColorOrange).Render("🔧 FIX SUGGESTION")
		fmt.Println(header)
		fmt.Println(ui.BodyStyleDouble(ui.ColorOrange).Render(response))

		// Offer to copy to clipboard
		fmt.Print("\n📋 Copy suggestion to clipboard? [y/N]: ")
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(strings.ToLower(input))
		if input == "y" || input == "yes" {
			if err := clipboard.WriteAll(response); err != nil {
				fmt.Println(ui.ErrorStyle().Render("❌ Could not copy to clipboard: " + err.Error()))
			} else {
				fmt.Println("✅ Copied to clipboard!")
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(fixCmd)
}
