package source

import (
	"crypto/sha256"
	"encoding/hex"
	"net/url"
)

type Source struct {
	Name        string
	Path        string
	DatabaseURL *url.URL
}

func NewSource(name, path, database string) (Source, error) {
	dbURL, err := url.Parse(database)
	if err != nil {
		return Source{}, err
	}

	return Source{
		Name:        name,
		Path:        path,
		DatabaseURL: dbURL,
	}, nil
}

func (s Source) Hash() (string, error) {
	hasher := sha256.New()
	_, err := hasher.Write([]byte(s.Path + s.DatabaseURL.String()))
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(hasher.Sum(nil)), nil
}
