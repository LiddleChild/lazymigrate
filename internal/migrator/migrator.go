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
	client *client

	path           string
	currentVersion uint
	isDirty        bool
}

func Open(path string, database string, verbose bool) (*Migrator, error) {
	sourceURL := fmt.Sprintf("file://%s", path)
	client, err := newClient(sourceURL, database)
	if err != nil {
		return nil, err
	}

	client.Log = newMigrateLogger(verbose)

	return &Migrator{
		client:         client,
		path:           path,
		currentVersion: 0,
		isDirty:        false,
	}, nil
}

func (m *Migrator) Stop() {
	m.client.GracefulStop <- true
}

func (m *Migrator) GetMigration() (*Migration, error) {
	currentVersion, isDirty, err := m.GetMigrationState()
	if err != nil {
		return nil, err
	}

	steps, err := m.fetchMigrations()
	if err != nil {
		return nil, err
	}

	m.currentVersion = currentVersion
	m.isDirty = isDirty

	return &Migration{
		Steps:          steps,
		CurrentVersion: currentVersion,
		IsDirty:        isDirty,
	}, nil
}

func (m *Migrator) GetMigrationState() (uint, bool, error) {
	currentVersion, isDirty, err := m.client.Version()
	switch {
	case errors.Is(err, migrate.ErrNilVersion):
		currentVersion = 0
		isDirty = false
	case err != nil:
		return 0, false, m.handleError(err)
	}

	return currentVersion, isDirty, nil
}

func (m *Migrator) MigrateToVersion(version uint) error {
	return m.handleError(m.client.Steps(int(version) - int(m.currentVersion)))
}

func (m *Migrator) ForceMigrateToVersion(version uint) error {
	return m.handleError(m.client.Force(int(version)))
}

func (m *Migrator) fetchMigrations() ([]MigrationStep, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	f, err := os.Open(m.path)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = f.Close()
	}()

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
			Path:     gopath.Join(wd, m.path, migration.Raw),
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

	return steps, nil
}

func (m *Migrator) handleError(err error) error {
	if err := m.client.Reconnect(); err != nil {
		return err
	}

	return err
}
