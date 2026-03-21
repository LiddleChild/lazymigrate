package sourceview

import (
	"charm.land/bubbles/v2/spinner"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/LiddleChild/lazymigrate/internal/appevent"
	"github.com/LiddleChild/lazymigrate/internal/brownsugar"
	"github.com/LiddleChild/lazymigrate/internal/source"
)

type Model struct {
	source           source.Source
	isLoadingContent bool

	spinner  spinner.Model
	viewport viewport.Model
}

func New() *Model {
	s := spinner.New()
	s.Spinner = spinner.MiniDot

	viewport := viewport.New()

	return &Model{
		source:           source.Source{},
		isLoadingContent: true,
		spinner:          s,
		viewport:         viewport,
	}
}

func (m *Model) Init() tea.Cmd {
	return tea.Batch(
		brownsugar.Cmd(appevent.NewUpdateSourcesRequestMsg()),
		m.spinner.Tick,
	)
}

func (m *Model) Update(msg tea.Msg) (*Model, tea.Cmd) {
	switch msg := msg.(type) {
	case appevent.UpdateSourcesMsg:
		m.source = msg.Sources[msg.CurrentSourceIndex]
		m.isLoadingContent = false

		return m, nil

	case spinner.TickMsg:
		if !m.isLoadingContent {
			return m, nil
		}

		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)

		return m, cmd
	}

	var cmd tea.Cmd
	m.spinner, cmd = m.spinner.Update(msg)

	return m, cmd
}

func (m *Model) Render(ctx brownsugar.Context) string {
	var (
		border = lipgloss.NewStyle().
			Width(ctx.Width).
			Height(ctx.Height).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.White)

		style = lipgloss.NewStyle()

		width  = ctx.Width - border.GetHorizontalFrameSize()
		height = ctx.Height - border.GetVerticalFrameSize()
	)

	if m.isLoadingContent {
		return border.
			AlignVertical(lipgloss.Center).
			AlignHorizontal(lipgloss.Center).
			Render(m.spinner.View())
	}

	m.viewport.SetWidth(width)
	m.viewport.SetHeight(height)
	m.viewport.SetContent(
		lipgloss.JoinVertical(lipgloss.Top,
			m.renderLine(style.Bold(true), m.source.Name),
			m.renderLine(style.Foreground(brownsugar.ColorBrightWhite), m.source.Path),
			m.renderLine(style.Foreground(brownsugar.ColorBrightWhite), m.source.DatabaseURL.Redacted()),
		),
	)

	return border.Render(m.viewport.View())
}

func (m *Model) renderLine(style lipgloss.Style, line string) string {
	return style.
		Width(lipgloss.Width(line)).
		Render(line)
}
