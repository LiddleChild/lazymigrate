package migrationview

import (
	"github.com/LiddleChild/lazymigrate/internal/migrator"
	tea "github.com/charmbracelet/bubbletea"
)

type UpdateMigrationMsg struct {
	Migration migrator.Migration
}

func UpdateMigrationsCmd(migration migrator.Migration) tea.Cmd {
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
