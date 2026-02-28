package contentview

import (
	"github.com/LiddleChild/lazymigrate/internal/migrator"
)

type updateMigrationContentMsg struct {
	MigrationStep migrator.MigrationStep
}
