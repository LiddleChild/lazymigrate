package source

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	gopath "path"
	"path/filepath"
	"strconv"

	"github.com/LiddleChild/lazymigrate/internal/cache"
	"github.com/LiddleChild/lazymigrate/internal/validator"
	"github.com/goccy/go-yaml"
)

type source struct {
	Name     string `yaml:"name"     validate:"required"`
	Path     string `yaml:"path"     validate:"required"`
	Database string `yaml:"database" validate:"required"`
}

type connectionFile struct {
	Sources []source `yaml:"sources" validate:"gt=0,dive,required"`
}

type Manager struct {
	cache *cache.Cache
	key   string

	index   int
	sources []Source
}

func NewManagerFromSource(cache *cache.Cache, path string, database string) (*Manager, error) {
	source, err := NewSource(gopath.Base(path), path, database)
	if err != nil {
		return nil, err
	}

	key, err := toCacheKey(path)
	if err != nil {
		return nil, err
	}

	return &Manager{
		cache:   cache,
		key:     key,
		index:   0,
		sources: []Source{source},
	}, nil
}

func NewManagerFromPath(cache *cache.Cache, path string) (*Manager, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = f.Close()
	}()

	var connFile connectionFile
	if err := yaml.NewDecoder(f).Decode(&connFile); err != nil {
		return nil, err
	}

	if err := validator.ValidateStruct(connFile); err != nil {
		return nil, err
	}

	sources := make([]Source, 0, len(connFile.Sources))
	for _, source := range connFile.Sources {
		s, err := NewSource(
			source.Name,
			source.Path,
			source.Database,
		)
		if err != nil {
			return nil, err
		}

		sources = append(sources, s)
	}

	key, err := toCacheKey(path)
	if err != nil {
		return nil, err
	}

	index, err := readIndexFromCache(cache, key)
	if err != nil {
		// ignore error and defaults index to 0
		index = 0
	}

	return &Manager{
		cache:   cache,
		key:     key,
		index:   index,
		sources: sources,
	}, nil
}

func (m *Manager) GetCurrentSource() Source {
	return m.sources[m.index]
}

func (m *Manager) GetCurrentSourceIndex() int {
	return m.index
}

func (m *Manager) SetCurrentSource(source Source) {
	for i, s := range m.sources {
		if s.Name == source.Name {
			m.index = i

			// ignore error
			_ = writeIndexToCache(m.cache, m.key, m.index)

			return
		}
	}
}

func (m *Manager) ListSources() []Source {
	return m.sources
}

func readIndexFromCache(cache *cache.Cache, key string) (int, error) {
	rawIndex, err := cache.Read(key)
	if err != nil {
		return 0, err
	}

	index, err := strconv.ParseInt(string(rawIndex), 10, 32)
	if err != nil {
		return 0, err
	}

	return int(index), nil
}

func writeIndexToCache(cache *cache.Cache, key string, index int) error {
	return cache.Write(key, []byte(strconv.FormatInt(int64(index), 10)))
}

func toCacheKey(path string) (string, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}

	hasher := sha256.New()
	if _, err := hasher.Write([]byte(absPath)); err != nil {
		return "", err
	}

	return hex.EncodeToString(hasher.Sum(nil)), nil
}
