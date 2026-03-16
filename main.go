package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/LiddleChild/lazymigrate/internal/app"
	"github.com/LiddleChild/lazymigrate/internal/appconfig"
	"github.com/LiddleChild/lazymigrate/internal/appevent"
	"github.com/LiddleChild/lazymigrate/internal/cache"
	"github.com/LiddleChild/lazymigrate/internal/log"
	"github.com/LiddleChild/lazymigrate/internal/migrator"
	"github.com/LiddleChild/lazymigrate/internal/runconfig"
	"github.com/LiddleChild/lazymigrate/internal/source"
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

	if cfg.Version {
		fmt.Println(appconfig.Name, appconfig.Version)
		return nil
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

	var sourcesManger *source.Manager
	if cfg.SourceFilePath != "" {
		sourcesManger, err = source.NewManagerFromPath(cache, cfg.SourceFilePath)
	} else {
		sourcesManger, err = source.NewManagerFromSource(cache, cfg.Path, cfg.Database)
	}
	if err != nil {
		return err
	}

	migrator := migrator.New(cache, cfg.IsVerbose)
	if err := migrator.Open(sourcesManger.GetCurrentSource()); err != nil {
		return err
	}

	app := app.New(migrator, sourcesManger)
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
