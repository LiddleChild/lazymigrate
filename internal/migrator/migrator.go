package migrator

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"maps"
	"os"
	gopath "path"
	"slices"
	"sync"

	"github.com/LiddleChild/lazymigrate/internal/cache"
	"github.com/LiddleChild/lazymigrate/internal/source"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	migratesource "github.com/golang-migrate/migrate/v4/source"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type migrator struct {
	client *client
	source source.Source

	currentVersion   uint
	isDirty          bool
	steps            []MigrationStep
	appliedMigration map[Signature]MigrationStep
}

type Migrator struct {
	*migrator
	sync.RWMutex

	cache   *cache.Cache
	verbose bool
}

func New(cache *cache.Cache, verbose bool) *Migrator {
	return &Migrator{
		cache:    cache,
		migrator: nil,
		verbose:  verbose,
	}
}

func (m *Migrator) Open(source source.Source) error {
	m.Lock()
	defer m.Unlock()

	if m.migrator != nil {
		m.client.Disconnect()
	}

	var (
		sourceURL   = fmt.Sprintf("file://%s", source.FullPath)
		databaseURL = source.DatabaseURL.String()
	)

	client, err := newClient(sourceURL, databaseURL, m.verbose)
	if err != nil {
		return err
	}

	currentVersion, isDirty, err := client.Version()
	switch {
	case errors.Is(err, migrate.ErrNilVersion):
		currentVersion = 0
		isDirty = false
	case err != nil:
		return err
	}

	steps, err := loadMigrations(source.FullPath)
	if err != nil {
		return err
	}

	cacheKey, err := source.Hash()
	if err != nil {
		return err
	}

	buffer, err := m.cache.Read(cacheKey)
	if err != nil {
		return err
	}

	appliedMigration := make(map[Signature]MigrationStep)
	if buffer != nil {
		if err := json.Unmarshal(buffer, &appliedMigration); err != nil {
			slog.Error(fmt.Sprintf("dirty cache state: %s", err.Error()))
		}
	}

	if len(appliedMigration) == 0 {
		appliedMigration = updateAppliedMigration(appliedMigration, updateAppliedMigrationParam{
			steps:       steps,
			fromVersion: 0,
			toVersion:   currentVersion,
			isDirty:     isDirty,
		})
	}

	m.migrator = &migrator{
		client:           client,
		source:           source,
		currentVersion:   currentVersion,
		isDirty:          isDirty,
		steps:            steps,
		appliedMigration: appliedMigration,
	}

	return nil
}

func (m *Migrator) Stop() {
	m.client.Disconnect()
}

func (m *Migrator) GetMigration() (*Migration, error) {
	var err error
	m.currentVersion, m.isDirty, err = m.GetMigrationState()
	if err != nil {
		return nil, err
	}

	m.steps, err = loadMigrations(m.source.FullPath)
	if err != nil {
		return nil, err
	}

	return &Migration{
		Steps:            m.steps,
		AppliedMigration: m.appliedMigration,
		CurrentVersion:   m.currentVersion,
		IsDirty:          m.isDirty,
	}, nil
}

func (m *Migrator) GetMigrationState() (uint, bool, error) {
	m.RLock()
	defer m.RUnlock()

	currentVersion, isDirty, err := m.client.Version()
	switch {
	case errors.Is(err, migrate.ErrNilVersion):
		currentVersion = 0
		isDirty = false
	case err != nil:
		return 0, false, m.handleError(err)
	}

	slog.Info("Fetched current migration state")

	return currentVersion, isDirty, nil
}

func (m *Migrator) MigrateToVersion(version uint) error {
	m.RLock()
	defer m.RUnlock()

	err := m.client.Steps(int(version) - int(m.currentVersion))
	if errors.Is(err, migrate.ErrNoChange) {
		return nil
	} else if err != nil {
		err := m.handleError(err)
		return errors.Join(err, m.updateAppliedMigration())
	}

	if err := m.updateAppliedMigration(); err != nil {
		return err
	}

	slog.Info(fmt.Sprintf("Migrated to version %d", version))

	return nil
}

func (m *Migrator) ForceMigrateToVersion(version uint) error {
	m.RLock()
	defer m.RUnlock()

	if err := m.client.Force(int(version)); err != nil {
		return m.handleError(err)
	}

	slog.Info(fmt.Sprintf("Forced to version %d", version))

	return nil
}

func (m *Migrator) CreateMigration(name string) error {
	var version uint = 1
	if len(m.steps) > 0 {
		version = m.steps[len(m.steps)-1].Version
		version++
	}

	for _, direction := range []string{"up", "down"} {
		filename := fmt.Sprintf("%06d_%s.%s.sql", version, name, direction)

		// same to migrate
		f, err := os.OpenFile(gopath.Join(m.source.FullPath, filename), os.O_RDWR|os.O_CREATE|os.O_EXCL, 0666)
		if err != nil {
			return err
		}

		if err := f.Close(); err != nil {
			return err
		}
	}

	slog.Info(fmt.Sprintf(`Created migration version %d "%s"`, version, name))

	return nil
}

func loadMigrations(path string) ([]MigrationStep, error) {
	f, err := os.Open(path)
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
		migration, err := migratesource.Parse(entry.Name())
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
				Signature:  NewSignature(),
			}
		}

		absPath := gopath.Join(path, migration.Raw)
		stepDirection := &MigrationStepDirection{
			Fullname: migration.Raw,
			Path:     absPath,
		}

		signature, err := NewSignatureFromFile(absPath)
		if err != nil {
			return nil, err
		}

		switch migration.Direction {
		case migratesource.Up:
			step.Up = stepDirection
			step.Signature = signature
		case migratesource.Down:
			step.Down = stepDirection
		}

		mp[migration.Version] = step
	}

	steps := []MigrationStep{}
	for _, version := range slices.Sorted(maps.Keys(mp)) {
		steps = append(steps, mp[version])
	}

	return steps, nil
}

func (m *Migrator) handleError(err error) error {
	if err != nil {
		if err := m.client.Reconnect(); err != nil {
			return err
		}
	}

	return err
}

func (m *Migrator) updateAppliedMigration() error {
	currentVersion, isDirty, err := m.GetMigrationState()
	if err != nil {
		return err
	}

	m.appliedMigration = updateAppliedMigration(m.appliedMigration, updateAppliedMigrationParam{
		steps:       m.steps,
		fromVersion: m.currentVersion,
		toVersion:   currentVersion,
		isDirty:     isDirty,
	})

	buffer, err := json.Marshal(m.appliedMigration)
	if err != nil {
		return err
	}

	cacheKey, err := m.source.Hash()
	if err != nil {
		return err
	}

	return m.cache.Write(cacheKey, buffer)
}

type updateAppliedMigrationParam struct {
	steps       []MigrationStep
	fromVersion uint
	toVersion   uint
	isDirty     bool
}

func updateAppliedMigration(appliedMigration map[Signature]MigrationStep, param updateAppliedMigrationParam) map[Signature]MigrationStep {
	if param.toVersion > param.fromVersion {
		// migrate up
		for _, step := range param.steps {
			if (step.Version > param.fromVersion && step.Version < param.toVersion) ||
				(!param.isDirty && step.Version == param.toVersion) {
				appliedMigration[step.Signature] = step
			}
		}
	} else {
		// migrate down
		for _, step := range param.steps {
			if step.Version > param.toVersion && step.Version <= param.fromVersion {
				delete(appliedMigration, step.Signature)
			}
		}
	}

	return appliedMigration
}
