package utils

import (
	"os"
	"path/filepath"
)

func IsFileNotExist(path string) bool {
	_, err := os.Stat(path)
	return err != nil && os.IsNotExist(err)
}

func CreateAndGetFile(path string) (*os.File, error) {
	return os.OpenFile(path, os.O_RDONLY|os.O_CREATE, 0644)
}

func CreateFile(path string) error {
	_, err := os.OpenFile(path, os.O_RDONLY|os.O_CREATE, 0644)
	return err
}

func Clone(pathList []string) error {
	for _, path := range pathList {
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			return err
		}
	}
	return nil
}
