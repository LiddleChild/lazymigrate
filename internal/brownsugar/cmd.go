package brownsugar

import tea "charm.land/bubbletea/v2"

func Cmd(msg tea.Msg) tea.Cmd {
	return func() tea.Msg {
		return msg
	}
}

type CmdAggregator []tea.Cmd

func (a *CmdAggregator) Add(cmd tea.Cmd) *CmdAggregator {
	*a = append(*a, cmd)
	return a
}
