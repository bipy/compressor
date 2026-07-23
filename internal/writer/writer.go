package writer

import (
	"fmt"
	"os"
	"path/filepath"
)

type Writer struct{}

func New() *Writer {
	return &Writer{}
}

func (w *Writer) Write(path string, data []byte) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("writer: mkdir %s: %w", dir, err)
	}

	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, 0644); err != nil {
		return fmt.Errorf("writer: write temp %s: %w", tmp, err)
	}

	if err := os.Rename(tmp, path); err != nil {
		os.Remove(tmp)
		return fmt.Errorf("writer: rename %s -> %s: %w", tmp, path, err)
	}

	return nil
}
