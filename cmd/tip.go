package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/shell-sage/internal/logger"
	"github.com/shell-sage/internal/metrics"
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

		lang := "English"
		if LangFlag != "" {
			lang = LangFlag
		}
		prompt := fmt.Sprintf(
			"IMPORTANT: You MUST respond ONLY in %s. Do not use any other language.\nGive me ONE practical, specific terminal/shell tip that most developers don't know. Be concise, max 3 sentences. No intro text.",
			lang,
		)

		pipe, err := buildPipeline()
		if err != nil {
			elapsed := time.Since(start)
			logger.Log.WithError(err).Error("'tip' failed to build pipeline")
			metrics.Record("tip", elapsed, err.Error())
			fmt.Println(ui.ErrorStyle().Render("❌ " + err.Error()))
			return
		}

		sp := spinner.New("Fetching a tip from the sage...")
		sp.Start()
		firstToken := true

		borderColor := lipgloss.Color("#FFD700")
		header := ui.HeaderStyle("#FFD700").Render("💡 TERMINAL TIP")

		response, err := pipe.RunStream(prompt, "tip", func(token string) {
			if firstToken {
				sp.Stop()
				firstToken = false
				fmt.Println(header)
				fmt.Println(lipgloss.NewStyle().Foreground(borderColor).Render("╭" + strings.Repeat("─", 76) + "╮"))
				fmt.Print(lipgloss.NewStyle().Foreground(borderColor).Render("│") + "  ")
			}
			formatted := strings.ReplaceAll(token, "\n", "\n"+lipgloss.NewStyle().Foreground(borderColor).Render("│")+"  ")
			fmt.Print(formatted)
		})

		if firstToken {
			sp.Stop()
		}

		elapsed := time.Since(start)

		if err != nil {
			if !firstToken {
				fmt.Println()
			}
			logger.Log.WithError(err).Error("'tip' command failed")
			metrics.Record("tip", elapsed, err.Error())
			fmt.Println(ui.ErrorStyle().Render("❌ " + err.Error()))
			return
		}

		fmt.Println()
		fmt.Println(lipgloss.NewStyle().Foreground(borderColor).Render("╰" + strings.Repeat("─", 76) + "╯"))

		logger.Log.WithField("duration_ms", elapsed.Milliseconds()).Info("'tip' command completed")
		metrics.Record("tip", elapsed, "")
		_ = response
	},
}

func init() {
	rootCmd.AddCommand(tipCmd)
}
