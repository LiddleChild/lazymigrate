package newmigrationview

import (
	"strings"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/LiddleChild/lazymigrate/internal/appevent"
	"github.com/LiddleChild/lazymigrate/internal/appscene"
	"github.com/LiddleChild/lazymigrate/internal/brownsugar"
)

var (
	KeyEnter = key.NewBinding(key.WithKeys("enter"))
	KeyEsc   = key.NewBinding(key.WithKeys("esc"))
)

type Model struct {
	input textinput.Model
}

func New() *Model {
	input := textinput.New()
	input.Focus()

	return &Model{
		input: input,
	}
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (*Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, KeyEnter):
			value := m.input.Value()
			m.input.Reset()

			value = strings.ReplaceAll(strings.TrimSpace(strings.ToLower(value)), " ", "_")

			return m, brownsugar.Cmd(appevent.NewCreateMigrationMsg(value))

		case key.Matches(msg, KeyEsc):
			return m, brownsugar.Cmd(appevent.NewSwitchSceneMsg(appscene.SceneHome))
		}

	case appevent.MigrationCreatedMsg:
		return m, brownsugar.Cmd(appevent.NewSwitchSceneMsg(appscene.SceneHome))
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m *Model) Render(ctx brownsugar.RenderContext) string {
	var (
		border = lipgloss.NewStyle().Border(lipgloss.RoundedBorder())

		window = border.Width(min(36, ctx.Width))
	)

	m.input.SetWidth(window.GetWidth() - window.GetHorizontalFrameSize() - len(m.input.Prompt) - 1)

	return lipgloss.NewStyle().
		Width(ctx.Width).
		Height(ctx.Height).
		AlignHorizontal(lipgloss.Center).
		AlignVertical(lipgloss.Center).
		Render(
			window.Render(
				lipgloss.JoinVertical(lipgloss.Top,
					"Name",
					m.input.View(),
				),
			),
		)
}
