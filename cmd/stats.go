package cmd

import (
	"fmt"
	"sort"

	"github.com/charmbracelet/lipgloss"
	"github.com/shell-sage/internal/metrics"
	"github.com/shell-sage/internal/ui"
	"github.com/spf13/cobra"
)

var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Show usage statistics for all ssage commands",
	Run: func(cmd *cobra.Command, args []string) {
		store := metrics.Load()

		if len(store) == 0 {
			fmt.Println(ui.ErrorStyle().Render("⚠️  No stats yet. Run some commands first!"))
			return
		}

		// Sort command names for consistent output
		names := make([]string, 0, len(store))
		for k := range store {
			names = append(names, k)
		}
		sort.Strings(names)

		// Styles
		titleStyle := lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#00D7FF")).
			MarginBottom(1)

		labelStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#888888")).
			Width(16)

		valueStyle := lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#E0E0E0"))

		failStyle := lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FF5555"))

		successStyle := lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#39FF14"))

		divider := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#333333")).
			Render("────────────────────────────────────────")

		fmt.Println(titleStyle.Render("📊  Shell Sage — Usage Statistics"))
		fmt.Println(divider)

		for _, name := range names {
			stat := store[name]

			icon := map[string]string{
				"explain": "⚡",
				"fix":     "🔧",
				"analyze": "🧠",
				"tip":     "💡",
			}[name]
			if icon == "" {
				icon = "▸"
			}

			failRate := 0
			if stat.Runs > 0 {
				failRate = (stat.Failures * 100) / stat.Runs
			}

			failsRendered := successStyle.Render(fmt.Sprintf("%d", stat.Failures))
			if stat.Failures > 0 {
				failsRendered = failStyle.Render(fmt.Sprintf("%d", stat.Failures))
			}

			fmt.Printf("\n%s %s\n",
				lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FFFFFF")).Render(icon+" ssage "+name),
				lipgloss.NewStyle().Foreground(lipgloss.Color("#555555")).Render(stat.LastRun.Format("last run: Jan 2 15:04")),
			)
			fmt.Printf("  %s %s\n", labelStyle.Render("Runs:"), valueStyle.Render(fmt.Sprintf("%d", stat.Runs)))
			fmt.Printf("  %s %s  (%d%% failure rate)\n", labelStyle.Render("Failures:"), failsRendered, failRate)
			fmt.Printf("  %s %s\n", labelStyle.Render("Avg Duration:"), valueStyle.Render(fmt.Sprintf("%dms", stat.AvgTimeMs)))
			if stat.LastError != "" {
				fmt.Printf("  %s %s\n", labelStyle.Render("Last Error:"),
					failStyle.Render(truncate(stat.LastError, 60)))
			}
			fmt.Println(divider)
		}
	},
}

// truncate shortens a string to maxLen chars, adding "…" if needed.
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-1] + "…"
}

func init() {
	rootCmd.AddCommand(statsCmd)
}
