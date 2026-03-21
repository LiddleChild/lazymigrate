package cache

import (
	"errors"
	"os"
	"path"

	"github.com/LiddleChild/lazymigrate/internal/appconfig"
)

type Cache struct {
	path string
}

func New() (*Cache, error) {
	basepath := path.Join(appconfig.CacheDirectoryPath, appconfig.Name)
	if err := os.MkdirAll(basepath, os.ModePerm); err != nil {
		return nil, err
	}

	return &Cache{
		path: basepath,
	}, nil
}

func (c *Cache) entryPathFromKey(key string) string {
	return path.Join(c.path, key)
}

func (c *Cache) Read(key string) ([]byte, error) {
	buffer, err := os.ReadFile(c.entryPathFromKey(key))
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return buffer, nil
}

func (c *Cache) Write(key string, value []byte) error {
	f, err := os.OpenFile(c.entryPathFromKey(key), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0o644)
	if err != nil {
		return err
	}

	defer func() {
		_ = f.Close()
	}()

	_, err = f.Write(value)
	return err
}
