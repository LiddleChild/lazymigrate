package newmigrationscene

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

var _ brownsugar.SceneModel = (*Model)(nil)

type Model struct {
	input textinput.Model
}

func New() brownsugar.SceneModel {
	input := textinput.New()
	input.Focus()

	return &Model{
		input: input,
	}
}

func (m *Model) Scene() string {
	return appscene.SceneNewMigration
}

func (m *Model) Init() tea.Cmd {
	m.input.SetValue("")
	return nil
}

func (m *Model) Update(msg tea.Msg) (brownsugar.SceneModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, KeyEnter):
			value := m.input.Value()
			m.input.Reset()

			value = strings.ReplaceAll(strings.TrimSpace(strings.ToLower(value)), " ", "_")

			return m, brownsugar.Cmd(appevent.NewCreateMigrationMsg(value))

		case key.Matches(msg, KeyEsc):
			return m, brownsugar.Cmd(brownsugar.NewSwitchSceneMsg(appscene.SceneHome))
		}

	case appevent.MigrationCreatedMsg:
		return m, brownsugar.Cmd(brownsugar.NewSwitchSceneMsg(appscene.SceneHome))
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m *Model) Render(ctx brownsugar.Context) string {
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
