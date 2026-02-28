package main

import (
	"fmt"
	"os"

	"github.com/LiddleChild/lazymigrate/internal/app"
	"github.com/LiddleChild/lazymigrate/internal/migrator"
	"github.com/LiddleChild/lazymigrate/internal/runconfig"
	"github.com/LiddleChild/lazymigrate/internal/validator"

	tea "github.com/charmbracelet/bubbletea"
)

// func main() {
// 	m, err := migrate.New(
// 		"file://test/migrations",
// 		"postgres://postgres:password@localhost:5432/postgres?sslmode=disable&x-migrations-table=merchant_schema_migrations")
// 	if err != nil {
// 		panic(err)
// 	}
//
// 	// err = m.Up()
// 	// if err != nil {
// 	// 	panic(err)
// 	// }
//
// 	version, dirty, err := m.Version()
// 	if err != nil {
// 		panic(err)
// 	}
//
// 	fmt.Println(version, dirty)
//
// 	fmt.Println("done")
// }

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

	p := tea.NewProgram(app, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		return err
	}

	return nil
}
