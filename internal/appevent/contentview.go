package appevent

import (
	tea "charm.land/bubbletea/v2"
	"github.com/LiddleChild/lazymigrate/internal/migrator"
)

type UpdateMigrationContentMsg struct {
	MigrationStep migrator.MigrationStep
}

func NewUpdateMigrationContentMsg(migrationStep migrator.MigrationStep) tea.Msg {
	return UpdateMigrationContentMsg{
		MigrationStep: migrationStep,
	}
}
