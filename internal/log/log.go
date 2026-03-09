package log

import (
	"io"
	"os"
	"path"

	"github.com/LiddleChild/lazymigrate/internal/appconfig"
)

var (
	Entry io.Writer
)

func Initialize(isDebug bool) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	basepath := path.Join(home, ".local/share", appconfig.Name)
	if err := os.MkdirAll(basepath, os.ModePerm); err != nil {
		return err
	}

	if isDebug {
		Entry, err = os.OpenFile(path.Join(basepath, "debug.log"), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o644)
		if err != nil {
			return err
		}
	} else {
		Entry, err = os.OpenFile(path.Join(basepath, "error.log"), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0o644)
		if err != nil {
			return err
		}
	}

	return nil
}
