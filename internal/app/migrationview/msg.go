package migrationview

import (
	tea "charm.land/bubbletea/v2"
	"github.com/LiddleChild/lazymigrate/internal/migrator"
)

type UpdateMigrationMsg struct {
	Migration *migrator.Migration
}

func UpdateMigrationCmd(migration *migrator.Migration) tea.Cmd {
	return func() tea.Msg {
		return UpdateMigrationMsg{
			Migration: migration,
		}
	}
}

type SelectMigrationStepMsg struct {
	MigrationStep migrator.MigrationStep
}

func SelectMigrationStepCmd(m migrator.MigrationStep) tea.Cmd {
	return func() tea.Msg {
		return SelectMigrationStepMsg{
			MigrationStep: m,
		}
	}
}

type MigrateMsg struct {
	Version uint
}

func MigrateCmd(version uint) tea.Cmd {
	return func() tea.Msg {
		return MigrateMsg{
			Version: version,
		}
	}
}

type ForceMigrateMsg struct {
	Version uint
}

func ForceMigrateCmd(version uint) tea.Cmd {
	return func() tea.Msg {
		return ForceMigrateMsg{
			Version: version,
		}
	}
}
