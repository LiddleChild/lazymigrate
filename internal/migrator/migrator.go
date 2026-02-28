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
	client *migrate.Migrate

	steps          []MigrationStep
	currentVersion uint
	isDirty        bool
}

func Open(path, database string) (*Migrator, error) {
	sourceURL := fmt.Sprintf("file://%s", path)
	client, err := migrate.New(sourceURL, database)
	if err != nil {
		return nil, err
	}

	var (
		currentVersion uint
		isDirty        bool
	)

	currentVersion, isDirty, err = client.Version()
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

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	entries, err := f.ReadDir(0)
	if err != nil {
		return nil, err
	}

	mp := make(map[uint]MigrationStep)
	for _, entry := range entries {
		m, err := source.Parse(entry.Name())
		if err != nil {
			// ignore that entry, if it cannot be parse
			// same behavior with migrate library
			continue
		}

		step, ok := mp[m.Version]
		if !ok {
			step = MigrationStep{
				Version:    m.Version,
				Identifier: m.Identifier,
			}
		}

		stepDirection := &MigrationStepDirection{
			Fullname: m.Raw,
			Path:     gopath.Join(wd, path, m.Raw),
		}

		switch m.Direction {
		case source.Up:
			step.Up = stepDirection
		case source.Down:
			step.Down = stepDirection
		}

		mp[m.Version] = step
	}

	migrations := []MigrationStep{}
	for _, m := range mp {
		migrations = append(migrations, m)
	}

	slices.SortFunc(migrations, func(a, b MigrationStep) int {
		return int(a.Version) - int(b.Version)
	})

	return &Migrator{
		client:         client,
		steps:          migrations,
		currentVersion: currentVersion,
		isDirty:        isDirty,
	}, nil
}

func (m *Migrator) GetMigration() Migration {
	return Migration{
		Steps:          m.steps,
		CurrentVersion: m.currentVersion,
		IsDirty:        m.isDirty,
	}
}

func (m *Migrator) Stop() {
	m.client.GracefulStop <- true
}

func (m *Migrator) Migrate(version uint) error {
	return m.client.Migrate(version)
}

func (m *Migrator) ForceMigrate(version uint) error {
	return m.client.Force(int(version))
}
