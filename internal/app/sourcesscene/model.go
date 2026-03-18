package sourcesscene

import (
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"github.com/LiddleChild/lazymigrate/internal/appevent"
	"github.com/LiddleChild/lazymigrate/internal/appscene"
	"github.com/LiddleChild/lazymigrate/internal/brownsugar"
	"github.com/LiddleChild/lazymigrate/internal/components/list"
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
	return tea.Batch(
		brownsugar.Cmd(appevent.NewUpdateSourcesRequestMsg()),
		brownsugar.Cmd(appevent.NewUpdateHelpMenuKeysMsg(m.HelpMenuBindings())),
	)
}

func (m *Model) Update(msg tea.Msg) (brownsugar.SceneModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, KeyMap.Back):
			return m, brownsugar.Cmd(brownsugar.NewSwitchSceneMsg(appscene.SceneHome))

		case key.Matches(msg, KeyMap.Down):
			m.list.Down()
			return m, nil

		case key.Matches(msg, KeyMap.Up):
			m.list.Up()
			return m, nil

		case key.Matches(msg, KeyMap.Select):
			item := m.list.GetSelectedItem().(item)
			if !item.current {
				return m, brownsugar.Cmd(appevent.NewChangeMigratorSourceMsg(item.Source))
			}

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
	m.list.FocusAtCursor()

	return m.list.Render(brownsugar.Context{
		Width:  ctx.Width,
		Height: ctx.Height,
	})
}

func (m *Model) HelpMenuBindings() []key.Binding {
	return []key.Binding{KeyMap.Back, KeyMap.Up, KeyMap.Down, KeyMap.Select}
}
