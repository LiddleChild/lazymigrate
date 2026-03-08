package brownsugar

import (
	tea "charm.land/bubbletea/v2"
)

type SceneModel interface {
	Scene() string

	Init() tea.Cmd

	Update(msg tea.Msg) (SceneModel, tea.Cmd)

	Render(ctx Context) string
}
