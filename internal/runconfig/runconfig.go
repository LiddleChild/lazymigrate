package runconfig

import (
	"flag"

	"github.com/LiddleChild/lazymigrate/internal/validator"
)

var (
	pathFlag           string
	databaseFlag       string
	sourceFilePathFlag string

	verboseFlag bool
	debugFlag   bool
)

func init() {
	flag.StringVar(&pathFlag, "path", "", "path to migrations")
	flag.StringVar(&databaseFlag, "database", "", "database connection string")
	flag.StringVar(&sourceFilePathFlag, "source-file", "", "path to source file")

	flag.BoolVar(&verboseFlag, "verbose", false, "verbose logging")
	flag.BoolVar(&debugFlag, "debug", false, "enable debug mode")
}

type RunConfig struct {
	Path           string `validate:"required_without=SourceFilePath"`
	Database       string `validate:"required_without=SourceFilePath"`
	SourceFilePath string `validate:"required_without=Path Database"`

	IsVerbose bool
	IsDebug   bool
}

func Parse() (RunConfig, error) {
	flag.Parse()

	cfg := RunConfig{
		Path:           pathFlag,
		Database:       databaseFlag,
		IsVerbose:      verboseFlag,
		IsDebug:        debugFlag,
		SourceFilePath: sourceFilePathFlag,
	}

	if err := validator.ValidateStruct(cfg); err != nil {
		return RunConfig{}, err
	}

	return cfg, nil
}
