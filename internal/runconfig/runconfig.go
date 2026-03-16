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

	versionFlag bool
)

func init() {
	flag.StringVar(&pathFlag, "path", "", "path to migrations")
	flag.StringVar(&databaseFlag, "database", "", "database connection string")
	flag.StringVar(&sourceFilePathFlag, "source-file", "", "path to source file")

	flag.BoolVar(&verboseFlag, "verbose", false, "verbose logging")
	flag.BoolVar(&debugFlag, "debug", false, "enable debug mode")

	flag.BoolVar(&versionFlag, "version", false, "show version")
}

type RunConfig struct {
	Path           string `validate:"required_without_all=SourceFilePath Version"`
	Database       string `validate:"required_without_all=SourceFilePath Version"`
	SourceFilePath string `validate:"required_without_all=Path Database Version"`

	IsVerbose bool
	IsDebug   bool

	Version bool
}

func Parse() (RunConfig, error) {
	flag.Parse()

	cfg := RunConfig{
		Path:           pathFlag,
		Database:       databaseFlag,
		SourceFilePath: sourceFilePathFlag,
		IsVerbose:      verboseFlag,
		IsDebug:        debugFlag,
		Version:        versionFlag,
	}

	if err := validator.ValidateStruct(cfg); err != nil {
		return RunConfig{}, err
	}

	return cfg, nil
}
