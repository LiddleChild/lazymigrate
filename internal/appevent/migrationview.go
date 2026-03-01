package appevent

import (
	tea "charm.land/bubbletea/v2"
	"github.com/LiddleChild/lazymigrate/internal/migrator"
)

type UpdateMigrationMsg struct {
	Migration *migrator.Migration
}

func NewUpdateMigrationMsg(migration *migrator.Migration) tea.Msg {
	return UpdateMigrationMsg{
		Migration: migration,
	}
}

type SelectMigrationStepMsg struct {
	MigrationStep migrator.MigrationStep
}

func NewSelectMigrationStepMsg(m migrator.MigrationStep) tea.Msg {
	return SelectMigrationStepMsg{
		MigrationStep: m,
	}
}

type MigrateMsg struct {
	Version uint
}

func NewMigrateMsg(version uint) tea.Msg {
	return MigrateMsg{
		Version: version,
	}
}

type ForceMigrateMsg struct {
	Version uint
}

func NewForceMigrateMsg(version uint) tea.Msg {
	return ForceMigrateMsg{
		Version: version,
	}
}
