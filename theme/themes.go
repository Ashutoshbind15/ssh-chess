package theme

import (
	"github.com/charmbracelet/lipgloss"
)

type Theme struct {
	renderer *lipgloss.Renderer
	base     lipgloss.Style
}

func BasicTheme(renderer *lipgloss.Renderer) Theme {
	return Theme{
		renderer: renderer,
		base: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("63")).
			Background(lipgloss.Color("6")),
	}
}

func (t Theme) Base() lipgloss.Style {
	return t.base
}
