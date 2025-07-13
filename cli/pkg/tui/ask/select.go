package ask

import (
	"github.com/charmbracelet/huh"
	"github.com/nitrictech/nitric/cli/internal/style/colors"
)

func NewSelect[T comparable]() *huh.Select[T] {
	return huh.NewSelect[T]().
		WithTheme(&colors.Theme).(*huh.Select[T])
}
