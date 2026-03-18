package appevent

import (
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
)

type UpdateMigrationRequestMsg struct{}

func NewUpdateMigrationRequestMsg() tea.Msg {
	return UpdateMigrationRequestMsg{}
}

type UpdateSourcesRequestMsg struct{}

func NewUpdateSourcesRequestMsg() tea.Msg {
	return UpdateSourcesRequestMsg{}
}

type UpdateHelpMenuKeysMsg struct {
	Bindings []key.Binding
}

func NewUpdateHelpMenuKeysMsg(bindings []key.Binding) tea.Msg {
	return UpdateHelpMenuKeysMsg{
		Bindings: bindings,
	}
}
