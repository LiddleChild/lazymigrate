package migrator

import (
	"errors"

	"github.com/golang-migrate/migrate/v4"
)

type client struct {
	*migrate.Migrate

	sourceURL   string
	databaseURL string
}

func newClient(sourceURL, databaseURL string) (*client, error) {
	migrate, err := migrate.New(sourceURL, databaseURL)
	if err != nil {
		return nil, err
	}

	return &client{
		Migrate:     migrate,
		sourceURL:   sourceURL,
		databaseURL: databaseURL,
	}, nil
}

func (conn *client) Reconnect() error {
	sourceErr, databaseErr := conn.Close()
	if err := errors.Join(sourceErr, databaseErr); err != nil {
		return err
	}

	var err error
	conn.Migrate, err = migrate.New(conn.sourceURL, conn.databaseURL)
	return err
}
