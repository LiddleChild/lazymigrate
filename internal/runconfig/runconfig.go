package runconfig

import (
	"flag"

	"github.com/LiddleChild/lazymigrate/internal/validator"
)

var (
	pathFlag     string
	databaseFlag string
	verboseFlag  bool
	debugFlag    bool
)

func init() {
	flag.StringVar(&pathFlag, "path", "", "path to migrations")
	flag.StringVar(&databaseFlag, "database", "", "database connection string")
	flag.BoolVar(&verboseFlag, "verbose", false, "verbose logging")
	flag.BoolVar(&debugFlag, "debug", false, "enable debug mode")
}

type RunConfig struct {
	Path      string `validate:"required"`
	Database  string `validate:"required"`
	IsVerbose bool
	IsDebug   bool
}

func Parse() (RunConfig, error) {
	flag.Parse()

	cfg := RunConfig{
		Path:      pathFlag,
		Database:  databaseFlag,
		IsVerbose: verboseFlag,
		IsDebug:   debugFlag,
	}

	if err := validator.ValidateStruct(cfg); err != nil {
		return RunConfig{}, err
	}

	return cfg, nil
}
