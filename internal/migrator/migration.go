package migrator

type Migration struct {
	Steps          []MigrationStep
	CurrentVersion uint
	IsDirty        bool
}

type MigrationStep struct {
	Version    uint
	Identifier string

	Up   *MigrationStepDirection
	Down *MigrationStepDirection
}

type MigrationStepDirection struct {
	Fullname string
	Path     string
}
