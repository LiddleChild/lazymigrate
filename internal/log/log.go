package log

import (
	"io"
	"os"
)

var (
	Entry io.Writer
)

func Initialize(isDebug bool) error {
	if isDebug {
		var err error
		Entry, err = os.OpenFile("debug.log", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o644)
		if err != nil {
			return err
		}
	} else {
		Entry = io.Discard
	}

	return nil
}
