package utils

import (
	"crypto/sha256"
	"io"
	"os"
)

func CreateHash(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	h := sha256.New()
	n, err := io.Copy(h, file)
	if err != nil {
		return "", err
	}
	if n == 0 {
		return "", nil // Return empty hash for empty files
	}
	return string(h.Sum(nil)), nil
}
