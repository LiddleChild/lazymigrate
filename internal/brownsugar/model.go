package brownsugar

import tea "charm.land/bubbletea/v2"

type ViewModel interface {
	Init() tea.Cmd

	Update(msg tea.Msg) (ViewModel, tea.Cmd)

	Render(ctx Context) string
}
