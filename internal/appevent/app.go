package appevent

import (
	tea "charm.land/bubbletea/v2"
)

type UpdateMigrationRequestMsg struct{}

func NewUpdateMigrationRequestMsg() tea.Msg {
	return UpdateMigrationRequestMsg{}
}
