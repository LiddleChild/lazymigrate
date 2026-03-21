package log

import (
	"io"
	"log/slog"
	"os"
	"path"

	"github.com/LiddleChild/lazymigrate/internal/appconfig"
)

var (
	Entry io.Writer
)

type LogFileWriter struct {
	entry io.Writer
}

func NewLogFileWriter(isDebug bool) (*LogFileWriter, error) {
	basepath := path.Join(appconfig.TempDirectoryPath, appconfig.Name)
	if err := os.MkdirAll(basepath, os.ModePerm); err != nil {
		return nil, err
	}

	entryPath := path.Join(basepath, "error.log")
	if isDebug {
		entryPath = path.Join(basepath, "debug.log")
	}

	entry, err := os.OpenFile(entryPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o644)
	if err != nil {
		return nil, err
	}

	return &LogFileWriter{
		entry: entry,
	}, nil
}

func (l *LogFileWriter) Handle(level slog.Level) slog.Handler {
	opts := slog.HandlerOptions{
		Level: level,
		ReplaceAttr: func(groups []string, attr slog.Attr) slog.Attr {
			return attr
		},
	}

	return slog.NewTextHandler(l.entry, &opts)
}
