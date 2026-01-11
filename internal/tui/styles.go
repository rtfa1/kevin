package tui

import "github.com/charmbracelet/lipgloss"

var (
	// Colors
	ColorPrimary = lipgloss.Color("62")
	ColorDim     = lipgloss.Color("240")
	ColorFocus   = lipgloss.Color("62")
	ColorBorder  = lipgloss.Color("238")

	// Base Styles
	ColumnStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorBorder).
			Padding(0, 1)

	FocusedColumnStyle = ColumnStyle.
				BorderForeground(ColorFocus)

	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorPrimary).
			Padding(0, 1)

	TaskStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder(), false, false, false, true). // Left border only
			BorderForeground(ColorDim).
			Padding(0, 1).
			MarginBottom(1)

	SelectedTaskStyle = TaskStyle.
				BorderForeground(ColorFocus).
				Foreground(ColorPrimary)
)
