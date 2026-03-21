package appconfig

import (
	"os"
	"path"
)

const (
	Name    = "lazymigrate"
	Version = "v1.1.0"
)

var (
	CacheDirectoryPath = home(".cache")
	TempDirectoryPath  = home(".local/share")
)

func home(paths ...string) string {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	paths = append([]string{home}, paths...)
	return path.Join(paths...)
}
