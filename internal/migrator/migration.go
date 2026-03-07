package migrator

type Migration struct {
	Steps            []MigrationStep
	AppliedMigration map[Signature]MigrationStep
	CurrentVersion   uint
	IsDirty          bool
}

type MigrationStep struct {
	Version    uint
	Identifier string
	Signature  Signature

	Up   *MigrationStepDirection
	Down *MigrationStepDirection
}

type MigrationStepDirection struct {
	Fullname string
	Path     string
}
