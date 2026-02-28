package app

import tea "charm.land/bubbletea/v2"

type updateMigrationRequestMsg struct{}

func updateMigrationRequestCmd() tea.Msg {
	return updateMigrationRequestMsg{}
}
