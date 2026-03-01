package appevent

import tea "charm.land/bubbletea/v2"

type ErrMsg struct {
	Err error
}

func ErrCmd(err error) tea.Cmd {
	return func() tea.Msg {
		return &ErrMsg{Err: err}
	}
}
