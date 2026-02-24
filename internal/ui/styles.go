package ui

import "github.com/charmbracelet/lipgloss"

var (
	StatusBarStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("252")).
			Background(lipgloss.Color("235")).
			Padding(0, 1)

	SpinnerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("205"))

	TimerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")).
			Bold(true)

	DimStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240"))

	LabelStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("252"))

	InputStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("252"))

	FocusedInputStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("205"))

	PhaseWorkingStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("205"))

	PhaseBreakStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("114"))
)
