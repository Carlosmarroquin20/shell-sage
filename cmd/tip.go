package cmd

import (
	"fmt"
	"time"

	"github.com/shell-sage/internal/logger"
	"github.com/shell-sage/internal/metrics"
	"github.com/shell-sage/internal/ollama"
	"github.com/shell-sage/internal/spinner"
	"github.com/shell-sage/internal/ui"
	"github.com/spf13/cobra"
)

var tipCmd = &cobra.Command{
	Use:   "tip",
	Short: "Get a quick, useful terminal tip from the AI",
	Run: func(cmd *cobra.Command, args []string) {
		start := time.Now()

		logger.Log.Info("Starting 'tip' command")

		sp := spinner.New("Fetching a tip from the sage...")
		sp.Start()

		client := ollama.NewClient(ModelFlag)

		// Build the prompt, optionally in a specific language
		lang := "English"
		if LangFlag != "" {
			lang = LangFlag
		}
		prompt := fmt.Sprintf(
			"Give me ONE practical, specific terminal/shell tip that most developers don't know. Be concise, max 3 sentences. Answer in %s. No intro text.",
			lang,
		)

		response, err := client.Generate(prompt)
		sp.Stop()

		elapsed := time.Since(start)

		if err != nil {
			logger.Log.WithError(err).Error("'tip' command failed")
			metrics.Record("tip", elapsed, err.Error())
			fmt.Println(ui.ErrorStyle().Render("❌ " + err.Error()))
			return
		}

		logger.Log.WithField("duration_ms", elapsed.Milliseconds()).Info("'tip' command completed")
		metrics.Record("tip", elapsed, "")

		header := ui.HeaderStyle("#FFD700").Render("💡 TERMINAL TIP")
		fmt.Println(header)
		fmt.Println(ui.BodyStyle("#FFD700").Render(response))
	},
}

func init() {
	rootCmd.AddCommand(tipCmd)
}
