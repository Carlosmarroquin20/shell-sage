package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/atotto/clipboard"
	"github.com/charmbracelet/lipgloss"
	"github.com/shell-sage/internal/logger"
	"github.com/shell-sage/internal/metrics"
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

		lang := "English"
		if LangFlag != "" {
			lang = LangFlag
		}
		prompt := fmt.Sprintf(
			"IMPORTANT: You MUST respond ONLY in %s. Do not use any other language.\nExplain this shell command in max 3 bullet points. Be extremely concise, no intro, no extra text: '%s'",
			lang, commandToExplain,
		)

		pipe, err := buildPipeline()
		if err != nil {
			elapsed := time.Since(start)
			logger.Log.WithError(err).Error("'explain' failed to build pipeline")
			metrics.Record("explain", elapsed, err.Error())
			fmt.Println(ui.ErrorStyle().Render("❌ " + err.Error()))
			return
		}

		// Show spinner until first token arrives
		sp := spinner.New("Consulting the AI sage...")
		sp.Start()
		firstToken := true

		borderColor := lipgloss.Color(ui.ColorCyan)
		header := ui.HeaderStyle(ui.ColorCyan).Render("⚡ EXPLAIN › " + commandToExplain)

		response, err := pipe.RunStream(prompt, "explain", func(token string) {
			if firstToken {
				sp.Stop()
				firstToken = false
				// Print header and open border
				fmt.Println(header)
				fmt.Println(lipgloss.NewStyle().Foreground(borderColor).Render("╭" + strings.Repeat("─", 76) + "╮"))
				fmt.Print(lipgloss.NewStyle().Foreground(borderColor).Render("│") + "  ")
			}
			// Print each token, handle newlines to keep box formatting
			formatted := strings.ReplaceAll(token, "\n", "\n"+lipgloss.NewStyle().Foreground(borderColor).Render("│")+"  ")
			fmt.Print(formatted)
		})

		if firstToken {
			sp.Stop() // In case we never got a token
		}

		elapsed := time.Since(start)

		if err != nil {
			if !firstToken {
				fmt.Println() // Clean newline after partial output
			}
			logger.Log.WithError(err).WithField("duration_ms", elapsed.Milliseconds()).Error("'explain' command failed")
			metrics.Record("explain", elapsed, err.Error())
			fmt.Println(ui.ErrorStyle().Render("❌ " + err.Error()))
			return
		}

		// Close the box
		fmt.Println()
		fmt.Println(lipgloss.NewStyle().Foreground(borderColor).Render("╰" + strings.Repeat("─", 76) + "╯"))

		logger.Log.WithField("duration_ms", elapsed.Milliseconds()).Info("'explain' command completed")
		metrics.Record("explain", elapsed, "")

		if CopyFlag {
			if err := clipboard.WriteAll(response); err != nil {
				logger.Log.WithError(err).Warn("Failed to copy to clipboard")
				fmt.Println(ui.ErrorStyle().Render("\n❌ Could not copy: " + err.Error()))
			} else {
				fmt.Println("\n✅ Copied to clipboard!")
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(explainCmd)
}
