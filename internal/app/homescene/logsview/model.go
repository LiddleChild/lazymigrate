package logsview

import (
	"strings"

	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/LiddleChild/lazymigrate/internal/appevent"
	"github.com/LiddleChild/lazymigrate/internal/brownsugar"
	"github.com/LiddleChild/lazymigrate/internal/components/focus"
	"github.com/LiddleChild/lazymigrate/internal/components/scrollpane"
	"github.com/LiddleChild/lazymigrate/internal/log"
)

type Model struct {
	focus.Model

	messages []log.Message

	viewport viewport.Model
}

func New() *Model {
	viewport := viewport.New()

	return &Model{
		Model:    focus.New(),
		messages: make([]log.Message, 0),
		viewport: viewport,
	}
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (*Model, tea.Cmd) {
	var (
		cmd tea.Cmd
		agg brownsugar.CmdAggregator
	)

	switch msg := msg.(type) {
	case appevent.LogMessageMsg:
		m.messages = append(m.messages, msg)
	}

	if m.IsFocused() {
		m.viewport, cmd = m.viewport.Update(msg)
		agg.Add(cmd)
	}

	return m, cmd
}

func (m *Model) Render(ctx brownsugar.Context) string {
	var (
		scrollpane = scrollpane.New().
				Foreground(m.borderColor()).
				BorderStyle(lipgloss.RoundedBorder()).
				CursorStyle(lipgloss.OuterHalfBlockBorder())

		width  = ctx.Width - scrollpane.GetHorizontalBorderSize()
		height = ctx.Height - scrollpane.GetVerticalBorderSize()
	)

	msgs := []string{}
	for _, msg := range m.messages {
		msgs = append(msgs, m.renderMessage(msg))
	}

	wasAtBottom := m.viewport.AtBottom()

	m.viewport.SetWidth(width)
	m.viewport.SetHeight(height)
	m.viewport.SetContent(lipgloss.NewStyle().
		Width(width).
		Render(lipgloss.JoinVertical(lipgloss.Top, msgs...)),
	)

	if wasAtBottom {
		m.viewport.GotoBottom()
	}

	return scrollpane.
		SetWidth(ctx.Width).
		SetHeight(ctx.Height).
		SetTotalLine(m.viewport.TotalLineCount()).
		SetCurrentLine(m.viewport.YOffset()).
		Render(m.viewport.View())
}

func (m *Model) borderColor() lipgloss.ANSIColor {
	if m.IsFocused() {
		return brownsugar.ColorYellow
	} else {
		return brownsugar.ColorWhite
	}
}

func (m *Model) renderMessage(msg log.Message) string {
	sTime := lipgloss.NewStyle().
		Foreground(brownsugar.ColorBrightBlack).
		Render(msg.Time.Format("2006-01-02 15:04:05"))

	var sLevel string
	switch msg.Level {
	case log.LogLevelDebug:
		sLevel = lipgloss.NewStyle().
			Foreground(brownsugar.ColorBrightMagenta).
			Render("DEBUG")

	case log.LogLevelInfo:
		sLevel = lipgloss.NewStyle().
			Foreground(brownsugar.ColorGreen).
			Render("INFO")

	case log.LogLevelWarn:
		sLevel = lipgloss.NewStyle().
			Foreground(brownsugar.ColorYellow).
			Render("WARN")

	case log.LogLevelError:
		sLevel = lipgloss.NewStyle().
			Foreground(brownsugar.ColorRed).
			Bold(true).
			Render("ERROR")
	}

	var sMsg string
	if msg.Secondary {
		sMsg = lipgloss.NewStyle().
			Foreground(brownsugar.ColorWhite).
			Render(msg.Message)
	} else {
		sMsg = lipgloss.NewStyle().
			Foreground(brownsugar.ColorBrightWhite).
			Bold(true).
			Render(msg.Message)
	}

	return strings.Join([]string{
		sTime,
		sLevel,
		sMsg,
	}, " ")
}
