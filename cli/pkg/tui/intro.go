package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/nitrictech/nitric/cli/internal/style"
	"github.com/nitrictech/nitric/cli/internal/style/icons"
	"github.com/nitrictech/nitric/cli/internal/version"
)

var suga = style.Purple(icons.Lightning + " " + version.ProductName)
var ver = style.Gray(version.GetShortVersion())

func SugaIntro(elements ...string) string {
	var b strings.Builder

	b.WriteString(suga + " " + ver + "\n")

	for i, element := range elements {
		isKey := i%2 == 0
		if isKey {
			b.WriteString("   - " + element + ": ")
		} else {
			b.WriteString(element)
			if i != len(elements)-1 {
				b.WriteString("\n")
			}
		}
	}

	return lipgloss.NewStyle().Border(lipgloss.HiddenBorder()).Render(b.String())
}
