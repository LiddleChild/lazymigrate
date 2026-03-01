package migrator

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/golang-migrate/migrate/v4"
)

var _ migrate.Logger = (*migrateLogger)(nil)

type migrateLogger struct {
	verbose bool
}

func newMigrateLogger(verbose bool) *migrateLogger {
	return &migrateLogger{
		verbose: verbose,
	}
}

func (l *migrateLogger) Printf(format string, v ...any) {
	msg, ok := strings.CutPrefix(fmt.Sprintf(format, v...), "error:")

	level := slog.LevelInfo
	if ok {
		level = slog.LevelError
	}

	slog.Log(context.Background(), level, strings.TrimSpace(msg))
}

func (l *migrateLogger) Verbose() bool {
	return l.verbose
}
