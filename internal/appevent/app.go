package appevent

import tea "charm.land/bubbletea/v2"

type UpdateMigrationRequestMsg struct{}

func UpdateMigrationRequestCmd() tea.Msg {
	return UpdateMigrationRequestMsg{}
}
