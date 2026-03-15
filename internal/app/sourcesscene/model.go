package sourcesscene

import (
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/LiddleChild/lazymigrate/internal/appevent"
	"github.com/LiddleChild/lazymigrate/internal/appscene"
	"github.com/LiddleChild/lazymigrate/internal/brownsugar"
	"github.com/LiddleChild/lazymigrate/internal/components/list"
)

var (
	KeyEsc  = key.NewBinding(key.WithKeys("esc"))
	KeyDown = key.NewBinding(key.WithKeys("j", "down"))
	KeyUp   = key.NewBinding(key.WithKeys("k", "up"))
)

var _ brownsugar.SceneModel = (*Model)(nil)

type Model struct {
	list *list.Model
}

func New() brownsugar.SceneModel {
	list := list.New()
	list.SetGap(1)

	return &Model{
		list: list,
	}
}

func (m *Model) Scene() string {
	return appscene.SceneSources
}

func (m *Model) Init() tea.Cmd {
	return brownsugar.Cmd(appevent.NewUpdateSourcesRequestMsg())
}

func (m *Model) Update(msg tea.Msg) (brownsugar.SceneModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, KeyEsc):
			return m, brownsugar.Cmd(brownsugar.NewSwitchSceneMsg(appscene.SceneHome))

		case key.Matches(msg, KeyDown):
			m.list.Down()
			return m, nil

		case key.Matches(msg, KeyUp):
			m.list.Up()
			return m, nil
		}

	case appevent.UpdateSourcesMsg:
		items := make([]list.Item, 0, len(msg.Sources))
		for i, source := range msg.Sources {
			items = append(items, item{
				Source:  source,
				current: i == msg.CurrentSourceIndex,
			})
		}

		m.list.SetItems(items)
		m.list.SetCursor(msg.CurrentSourceIndex)
	}

	return m, nil
}

func (m *Model) Render(ctx brownsugar.Context) string {
	var (
		width  = min(64, ctx.Width)
		height = min(32, ctx.Height) // height is roughly 2 times bigger than width

		window = lipgloss.NewStyle().
			AlignHorizontal(lipgloss.Center).
			AlignVertical(lipgloss.Center).
			Width(ctx.Width).
			Height(ctx.Height)
	)

	m.list.FocusAtCursor()

	return window.Render(
		lipgloss.JoinVertical(lipgloss.Top,
			"Sources",
			m.list.Render(brownsugar.Context{
				Width:  width,
				Height: height,
			}),
		),
	)
}
