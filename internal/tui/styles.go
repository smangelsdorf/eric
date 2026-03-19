package tui

import "github.com/charmbracelet/lipgloss"

var (
	// Table
	tableBorderColor = lipgloss.Color("6") // cyan
	tableBorder      = lipgloss.NewStyle().
				BorderStyle(lipgloss.RoundedBorder()).
				BorderForeground(tableBorderColor)
	headerStyle    = lipgloss.NewStyle().Bold(true).Padding(0, 1)
	cellStyle      = lipgloss.NewStyle().Padding(0, 1)
	cursorStyle    = lipgloss.NewStyle().Reverse(true)
	closedRowStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("8")) // grey
	statusOpen     = lipgloss.NewStyle().Foreground(lipgloss.Color("2")) // green
	statusClosed   = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	dateOld        = lipgloss.NewStyle().Foreground(lipgloss.Color("3")) // yellow

	// Modal
	modalBorder = lipgloss.NewStyle().BorderStyle(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("7")).Padding(1, 2) // white

	// Hint bar
	hintKey   = lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Bold(true) // light green
	hintLabel = lipgloss.NewStyle().Foreground(lipgloss.Color("7"))            // white
)
