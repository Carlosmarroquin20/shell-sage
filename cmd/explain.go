package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/shell-sage/internal/logger"
	"github.com/shell-sage/internal/metrics"
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
		start := time.Now()
		commandToExplain := strings.Join(args, " ")

		logger.Log.WithField("command", commandToExplain).Info("Starting 'explain' command")

		sp := spinner.New("Consulting the AI sage...")
		sp.Start()

		client := ollama.NewClient(ModelFlag)
		prompt := fmt.Sprintf("Explain this shell command in max 3 bullet points. Be extremely concise, no intro, no extra text: '%s'", commandToExplain)

		response, err := client.Generate(prompt)
		sp.Stop()

		elapsed := time.Since(start)

		if err != nil {
			logger.Log.WithError(err).WithField("duration_ms", elapsed.Milliseconds()).Error("'explain' command failed")
			metrics.Record("explain", elapsed, err.Error())
			fmt.Println(ui.ErrorStyle().Render("❌ " + err.Error()))
			return
		}

		logger.Log.WithField("duration_ms", elapsed.Milliseconds()).Info("'explain' command completed successfully")
		metrics.Record("explain", elapsed, "")

		header := ui.HeaderStyle(ui.ColorCyan).Render("⚡ EXPLAIN › " + commandToExplain)
		fmt.Println(header)
		fmt.Println(ui.BodyStyle(ui.ColorCyan).Render(response))
	},
}

func init() {
	rootCmd.AddCommand(explainCmd)
}
