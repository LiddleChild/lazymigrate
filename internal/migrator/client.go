package migrator

import (
	"errors"
	"log/slog"

	"github.com/golang-migrate/migrate/v4"
)

type client struct {
	*migrate.Migrate

	sourceURL   string
	databaseURL string
	verbose     bool
}

func newClient(sourceURL, databaseURL string, verbose bool) (*client, error) {
	cli := &client{
		Migrate:     nil,
		sourceURL:   sourceURL,
		databaseURL: databaseURL,
		verbose:     verbose,
	}

	if err := cli.Connect(); err != nil {
		return nil, err
	}

	return cli, nil
}

func (conn *client) Reconnect() error {
	sourceErr, databaseErr := conn.Close()
	if err := errors.Join(sourceErr, databaseErr); err != nil {
		return err
	}

	slog.Info("Connection closed")

	return conn.Connect()
}

func (conn *client) Connect() error {
	var err error
	conn.Migrate, err = migrate.New(conn.sourceURL, conn.databaseURL)
	if err != nil {
		return err
	}

	conn.Log = newMigrateLogger(conn.verbose)

	slog.Info("Connected")

	return nil
}

func (conn *client) Disconnect() {
	conn.GracefulStop <- true
}
