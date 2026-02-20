package cmd

import (
	"fmt"
	"strings"

	"github.com/shell-sage/internal/ollama"
	"github.com/shell-sage/internal/spinner"
	"github.com/shell-sage/internal/ui"
	"github.com/spf13/cobra"
)

var explainCmd = &cobra.Command{
	Use:   "explain [command]",
	Short: "Explain a shell command",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		commandToExplain := strings.Join(args, " ")

		sp := spinner.New("Consulting the AI sage...")
		sp.Start()

		client := ollama.NewClient(ModelFlag)
		prompt := fmt.Sprintf("Explain this shell command in max 3 bullet points. Be extremely concise, no intro, no extra text: '%s'", commandToExplain)

		response, err := client.Generate(prompt)
		sp.Stop()

		if err != nil {
			fmt.Println(ui.ErrorStyle().Render("❌ " + err.Error()))
			return
		}

		header := ui.HeaderStyle(ui.ColorCyan).Render("⚡ EXPLAIN › " + commandToExplain)
		fmt.Println(header)
		fmt.Println(ui.BodyStyle(ui.ColorCyan).Render(response))
	},
}

func init() {
	rootCmd.AddCommand(explainCmd)
}
