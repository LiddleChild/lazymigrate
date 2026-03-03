package appevent

import tea "charm.land/bubbletea/v2"

type CreateMigrationMsg struct {
	Name string
}

func NewCreateMigrationMsg(name string) tea.Msg {
	return CreateMigrationMsg{
		Name: name,
	}
}

type MigrationCreatedMsg struct{}

func NewMigrationCreatedMsg() tea.Msg {
	return MigrationCreatedMsg{}
}
