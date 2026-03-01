package appevent

import tea "charm.land/bubbletea/v2"

type ErrMsg struct {
	Err error
}

func NewErrMsg(err error) tea.Msg {
	return ErrMsg{
		Err: err,
	}
}
