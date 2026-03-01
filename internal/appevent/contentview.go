package appevent

import (
	"github.com/LiddleChild/lazymigrate/internal/migrator"
)

type UpdateMigrationContentMsg struct {
	MigrationStep migrator.MigrationStep
}
