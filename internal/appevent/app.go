package appevent

import (
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
