package style

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/nitrictech/nitric/cli/internal/style/colors"
)

var (
	White  = lipgloss.NewStyle().Foreground(colors.White).Render
	Gray   = lipgloss.NewStyle().Foreground(colors.Gray).Render
	Black  = lipgloss.NewStyle().Foreground(colors.Black).Render
	Red    = lipgloss.NewStyle().Foreground(colors.Red).Render
	Orange = lipgloss.NewStyle().Foreground(colors.Orange).Render
	Yellow = lipgloss.NewStyle().Foreground(colors.Yellow).Render
	Green  = lipgloss.NewStyle().Foreground(colors.Green).Render
	Teal   = lipgloss.NewStyle().Foreground(colors.Teal).Render
	Blue   = lipgloss.NewStyle().Foreground(colors.Blue).Render
	Purple = lipgloss.NewStyle().Foreground(colors.Purple).Render
)
