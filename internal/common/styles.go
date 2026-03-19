package common

import "github.com/charmbracelet/lipgloss"

// Game colors
const (
	ColorGreen  = lipgloss.Color("#538D4E")
	ColorYellow = lipgloss.Color("#B59F3B")
	ColorGray   = lipgloss.Color("#3A3A3C")
	ColorWhite  = lipgloss.Color("#FFFFFF")
	ColorRed    = lipgloss.Color("#CC3333")
	ColorBlue   = lipgloss.Color("#4A90D9")
	ColorPurple = lipgloss.Color("#9B59B6")
)

// Guess game styles
var (
	CorrectStyle = lipgloss.NewStyle().
			Background(ColorGreen).
			Foreground(ColorWhite).
			Bold(true)

	PresentStyle = lipgloss.NewStyle().
			Background(ColorYellow).
			Foreground(ColorWhite).
			Bold(true)

	AbsentStyle = lipgloss.NewStyle().
			Background(ColorGray).
			Foreground(ColorWhite)

	EmptyStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(ColorGray)
)

// Group game styles — four difficulty tiers from easiest to hardest
var TierStyles = [4]lipgloss.Style{
	lipgloss.NewStyle().Background(ColorYellow).Foreground(ColorWhite).Bold(true),
	lipgloss.NewStyle().Background(ColorGreen).Foreground(ColorWhite).Bold(true),
	lipgloss.NewStyle().Background(ColorBlue).Foreground(ColorWhite).Bold(true),
	lipgloss.NewStyle().Background(ColorPurple).Foreground(ColorWhite).Bold(true),
}

var SelectedStyle = lipgloss.NewStyle().
	Border(lipgloss.ThickBorder()).
	BorderForeground(ColorWhite)

// Shared UI styles
var (
	ErrorStyle = lipgloss.NewStyle().Foreground(ColorRed)

	SuccessStyle = lipgloss.NewStyle().Foreground(ColorGreen)

	TitleStyle = lipgloss.NewStyle().Bold(true).MarginBottom(1)

	HelpStyle = lipgloss.NewStyle().Faint(true)
)

// CellStyle returns a base cell style with the given dimensions, centered and padded.
func CellStyle(width, height int) lipgloss.Style {
	return lipgloss.NewStyle().
		Width(width).
		Height(height).
		Align(lipgloss.Center, lipgloss.Center).
		PaddingLeft(1).
		PaddingRight(1)
}
