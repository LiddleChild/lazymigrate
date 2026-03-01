package logsview

import (
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/LiddleChild/lazymigrate/internal/brownsugar"
	"github.com/LiddleChild/lazymigrate/internal/components/focus"
)

type Model struct {
	focus.Model

	viewport viewport.Model
}

func New() *Model {
	viewport := viewport.New()

	return &Model{
		Model:    focus.New(),
		viewport: viewport,
	}
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (*Model, tea.Cmd) {
	return m, nil
}

func (m *Model) Render(ctx brownsugar.RenderContext) string {
	var (
		border = lipgloss.NewStyle().Border(lipgloss.RoundedBorder())

		width  = ctx.Width - border.GetHorizontalFrameSize()
		height = ctx.Height - border.GetVerticalFrameSize()
	)

	m.viewport.SetWidth(width)
	m.viewport.SetHeight(height)
	m.viewport.SetContent("logs")

	return border.
		BorderForeground(m.borderColor()).
		Render(m.viewport.View())
}

func (m *Model) borderColor() lipgloss.ANSIColor {
	if m.IsFocused() {
		return brownsugar.ColorYellow
	} else {
		return brownsugar.ColorWhite
	}
}
