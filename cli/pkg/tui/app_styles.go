package tui

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/nitrictech/nitric/cli/internal/style/colors"
)

type AppStyles struct {
	Emphasize lipgloss.Style
	Faint     lipgloss.Style
	Success   lipgloss.Style
}

func NewAppStyles() AppStyles {
	return AppStyles{
		Emphasize: lipgloss.NewStyle().Foreground(colors.Teal).Bold(true),
		Faint:     lipgloss.NewStyle().Faint(true),
		Success:   lipgloss.NewStyle().Foreground(colors.Teal).Bold(true),
	}
}
