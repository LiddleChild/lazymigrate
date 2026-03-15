package source

import (
	"crypto/sha256"
	"encoding/hex"
	"net/url"
	"path/filepath"
)

type Source struct {
	Name        string
	Path        string
	FullPath    string
	DatabaseURL *url.URL
}

func NewSource(name, path, database string) (Source, error) {
	dbURL, err := url.Parse(database)
	if err != nil {
		return Source{}, err
	}

	fullPath, err := filepath.Abs(path)
	if err != nil {
		return Source{}, err
	}

	return Source{
		Name:        name,
		Path:        path,
		FullPath:    fullPath,
		DatabaseURL: dbURL,
	}, nil
}

func (s Source) Hash() (string, error) {
	hasher := sha256.New()
	_, err := hasher.Write([]byte(s.FullPath + s.DatabaseURL.String()))
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(hasher.Sum(nil)), nil
}
