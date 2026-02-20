// Package ui provides shared visual styles for the shell-sage CLI using lipgloss.
// All commands should use these definitions instead of defining their own.
package ui

import "github.com/charmbracelet/lipgloss"

// Color palette
const (
	ColorCyan   = "#00D7FF"
	ColorOrange = "#FF8C00"
	ColorGreen  = "#39FF14"
	ColorDim    = "#1a1a2e"
	ColorText   = "#E0E0E0"
)

// HeaderStyle returns a styled header label with the given accent color.
func HeaderStyle(color string) lipgloss.Style {
	return lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(color)).
		Background(lipgloss.Color(ColorDim)).
		Padding(0, 1)
}

// BodyStyle returns a bordered body box with the given accent color.
func BodyStyle(color string) lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorText)).
		PaddingTop(1).
		PaddingBottom(1).
		PaddingLeft(2).
		PaddingRight(2).
		Width(76).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(color))
}

// BodyStyleDouble returns a double-bordered body box — used for the fix command.
func BodyStyleDouble(color string) lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorText)).
		PaddingTop(1).
		PaddingBottom(1).
		PaddingLeft(2).
		PaddingRight(2).
		Width(76).
		BorderStyle(lipgloss.DoubleBorder()).
		BorderForeground(lipgloss.Color(color))
}

// ErrorStyle returns a style for error messages.
func ErrorStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FF5555"))
}
