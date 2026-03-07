package migrator

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"os"
	gopath "path"
)

type Signature string

func NewSignature() Signature {
	return ""
}

func NewSignatureFromFile(path string) (Signature, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}

	hasher := md5.New()
	if _, err := hasher.Write([]byte(gopath.Base(path))); err != nil {
		return "", err
	}

	if _, err := io.Copy(hasher, f); err != nil {
		return "", err
	}

	return Signature(hex.EncodeToString(hasher.Sum(nil))), nil
}
