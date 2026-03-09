package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/LiddleChild/lazymigrate/internal/app"
	"github.com/LiddleChild/lazymigrate/internal/appevent"
	"github.com/LiddleChild/lazymigrate/internal/cache"
	"github.com/LiddleChild/lazymigrate/internal/log"
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

	if err := log.Initialize(cfg.IsDebug); err != nil {
		return err
	}

	logDispatcher := log.NewLogDispatcher()

	var handlerOpt slog.HandlerOptions
	if cfg.IsDebug {
		handlerOpt.Level = slog.LevelDebug
	} else {
		handlerOpt.Level = slog.LevelError
	}

	slog.SetDefault(
		slog.New(slog.NewMultiHandler(
			slog.NewTextHandler(log.Entry, &handlerOpt),
			logDispatcher.Handler(),
		)),
	)

	cache, err := cache.New()
	if err != nil {
		return err
	}

	migrator, err := migrator.New(cache, cfg.Path, cfg.Database, cfg.IsVerbose)
	if err != nil {
		return err
	}

	app := app.New(migrator)
	p := tea.NewProgram(app)

	go func() {
		for msg := range logDispatcher.Pull() {
			p.Send(appevent.NewLogMessageMsg(msg))
		}
	}()

	if _, err := p.Run(); err != nil {
		return err
	}

	return nil
}
