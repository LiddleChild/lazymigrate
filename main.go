package main

import (
	"fmt"
	"os"

	"github.com/LiddleChild/lazymigrate/internal/app"
	"github.com/LiddleChild/lazymigrate/internal/migrator"
	"github.com/LiddleChild/lazymigrate/internal/runconfig"
	"github.com/LiddleChild/lazymigrate/internal/validator"

	tea "charm.land/bubbletea/v2"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func run() error {
	validator.Initialize()

	cfg, err := runconfig.Parse()
	if err != nil {
		return err
	}

	m, err := migrator.Open(cfg.Path, cfg.Database)
	if err != nil {
		return err
	}

	app := app.New(m)

	if _, err := tea.NewProgram(app).Run(); err != nil {
		return err
	}

	return nil
}
