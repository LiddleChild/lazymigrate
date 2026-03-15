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
	index   int
	sources []Source
}

func NewManagerFromSource(path string, database string) (*Manager, error) {
	source, err := NewSource(gopath.Base(path), path, database)
	if err != nil {
		return nil, err
	}

	return &Manager{
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

	index, err := readIndexFromCache(cache, path)
	if err != nil {
		// ignore error and defaults index to 0
		index = 0
	}

	return &Manager{
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

func (m *Manager) ListSources() []Source {
	return m.sources
}

func readIndexFromCache(cache *cache.Cache, path string) (int, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return 0, err
	}

	hasher := sha256.New()
	if _, err := hasher.Write([]byte(absPath)); err != nil {
		return 0, err
	}

	rawIndex, err := cache.Read(hex.EncodeToString(hasher.Sum(nil)))
	if err != nil {
		return 0, err
	}

	index, err := strconv.ParseInt(string(rawIndex), 10, 32)
	if err != nil {
		return 0, err
	}

	return int(index), nil
}
