package migrator

import (
	"errors"
	"fmt"
	"os"
	gopath "path"
	"slices"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type Migrator struct {
	client     *migrate.Migrate
	sourcePath string
	migration  *Migration
}

func Open(path, database string) (*Migrator, error) {
	sourceURL := fmt.Sprintf("file://%s", path)
	client, err := migrate.New(sourceURL, database)
	if err != nil {
		return nil, err
	}

	return &Migrator{
		sourcePath: path,
		client:     client,
	}, nil
}

func (m *Migrator) Stop() {
	m.client.GracefulStop <- true
}

func (m *Migrator) GetMigration() (*Migration, error) {
	migration, err := m.refreshMigration()
	if err != nil {
		return nil, err
	}

	m.migration = migration

	return migration, nil
}

func (m *Migrator) MigrateToVersion(version uint) error {
	return m.client.Steps(int(version) - int(m.migration.CurrentVersion))
}

func (m *Migrator) ForceMigrateToVersion(version uint) error {
	return m.client.Force(int(version))
}

func (m *Migrator) refreshMigration() (*Migration, error) {
	var (
		currentVersion uint
		isDirty        bool
	)

	currentVersion, isDirty, err := m.client.Version()
	switch {
	case errors.Is(err, migrate.ErrNilVersion):
		currentVersion = 0
		isDirty = false
	case err != nil:
		return nil, err
	}

	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	f, err := os.Open(m.sourcePath)
	if err != nil {
		return nil, err
	}

	entries, err := f.ReadDir(0)
	if err != nil {
		return nil, err
	}

	mp := make(map[uint]MigrationStep)
	for _, entry := range entries {
		migration, err := source.Parse(entry.Name())
		if err != nil {
			// ignore that entry, if it cannot be parse
			// same behavior with migrate library
			continue
		}

		step, ok := mp[migration.Version]
		if !ok {
			step = MigrationStep{
				Version:    migration.Version,
				Identifier: migration.Identifier,
			}
		}

		stepDirection := &MigrationStepDirection{
			Fullname: migration.Raw,
			Path:     gopath.Join(wd, m.sourcePath, migration.Raw),
		}

		switch migration.Direction {
		case source.Up:
			step.Up = stepDirection
		case source.Down:
			step.Down = stepDirection
		}

		mp[migration.Version] = step
	}

	steps := []MigrationStep{}
	for _, m := range mp {
		steps = append(steps, m)
	}

	slices.SortFunc(steps, func(a, b MigrationStep) int {
		return int(a.Version) - int(b.Version)
	})

	return &Migration{
		Steps:          steps,
		CurrentVersion: currentVersion,
		IsDirty:        isDirty,
	}, nil
}
