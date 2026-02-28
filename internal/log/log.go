package log

import "os"

var (
	Entry *os.File
)

func init() {
	var err error
	Entry, err = os.OpenFile("debug.log", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o644)
	if err != nil {
		panic(err)
	}
}
