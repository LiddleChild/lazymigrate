package migrator

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

func makeAppliedMigrations(steps []MigrationStep, indices ...int) map[Signature]MigrationStep {
	applied := make(map[Signature]MigrationStep)
	for _, step := range steps {
		if slices.Contains(indices, int(step.Version)) {
			applied[step.Signature] = step
		}
	}
	return applied
}

func Test_updateAppliedMigrations(t *testing.T) {
	steps := []MigrationStep{
		{
			Version:    1,
			Identifier: "migration_1",
			Signature:  Signature("1"),
			Up:         nil,
			Down:       nil,
		},
		{
			Version:    2,
			Identifier: "migration_2",
			Signature:  Signature("2"),
			Up:         nil,
			Down:       nil,
		},
		{
			Version:    3,
			Identifier: "migration_3",
			Signature:  Signature("3"),
			Up:         nil,
			Down:       nil,
		},
	}

	testcases := []struct {
		name                      string
		fromVersion               uint
		toVersion                 uint
		isDirty                   bool
		originalAppliedMigrations map[Signature]MigrationStep
		expectedAppliedMigrations map[Signature]MigrationStep
	}{
		{
			name:                      "successful up migration",
			fromVersion:               0,
			toVersion:                 3,
			isDirty:                   false,
			originalAppliedMigrations: makeAppliedMigrations(steps),
			expectedAppliedMigrations: makeAppliedMigrations(steps, 1, 2, 3),
		},
		{
			name:                      "dirty up migration",
			fromVersion:               0,
			toVersion:                 3,
			isDirty:                   true,
			originalAppliedMigrations: makeAppliedMigrations(steps),
			expectedAppliedMigrations: makeAppliedMigrations(steps, 1, 2),
		},
		{
			name:                      "successful down migration",
			fromVersion:               3,
			toVersion:                 0,
			isDirty:                   false,
			originalAppliedMigrations: makeAppliedMigrations(steps, 1, 2, 3),
			expectedAppliedMigrations: makeAppliedMigrations(steps),
		},
		{
			name:                      "dirty down migration",
			fromVersion:               3,
			toVersion:                 2,
			isDirty:                   true,
			originalAppliedMigrations: makeAppliedMigrations(steps, 1, 2, 3),
			expectedAppliedMigrations: makeAppliedMigrations(steps, 1, 2),
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			applied := updateAppliedMigration(tc.originalAppliedMigrations, updateAppliedMigrationParam{
				steps:       steps,
				fromVersion: tc.fromVersion,
				toVersion:   tc.toVersion,
				isDirty:     tc.isDirty,
			})
			assert.EqualValues(t, tc.expectedAppliedMigrations, applied)
		})
	}
}
