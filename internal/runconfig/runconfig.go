package runconfig

import (
	"flag"

	"github.com/LiddleChild/lazymigrate/internal/validator"
)

var (
	pathFlag     string
	databaseFlag string
)

func init() {
	flag.StringVar(&pathFlag, "path", "", "path to migrations")
	flag.StringVar(&databaseFlag, "database", "", "database connection string")
}

type RunConfig struct {
	Path     string `validate:"required"`
	Database string `validate:"required"`
}

func Parse() (RunConfig, error) {
	flag.Parse()

	cfg := RunConfig{
		Path:     pathFlag,
		Database: databaseFlag,
	}

	if err := validator.ValidateStruct(cfg); err != nil {
		return RunConfig{}, err
	}

	return cfg, nil
}
