package writer

import (
	"os"
	"path/filepath"
	"testing"
)

func TestWriter_Write(t *testing.T) {
	w := New()
	dir := t.TempDir()
	path := filepath.Join(dir, "sub", "test.txt")

	data := []byte("hello world")
	if err := w.Write(path, data); err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile failed: %v", err)
	}
	if string(got) != "hello world" {
		t.Fatalf("got %q, want %q", got, "hello world")
	}

	tmp := path + ".tmp"
	if _, err := os.Stat(tmp); !os.IsNotExist(err) {
		t.Fatalf("temp file should not exist after successful write")
	}
}

func TestWriter_WriteCreatesDirectories(t *testing.T) {
	w := New()
	dir := t.TempDir()
	path := filepath.Join(dir, "a", "b", "c", "file.txt")

	if err := w.Write(path, []byte("data")); err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	if _, err := os.Stat(path); err != nil {
		t.Fatalf("file should exist: %v", err)
	}
}

func TestWriter_WriteOverwritesExisting(t *testing.T) {
	w := New()
	dir := t.TempDir()
	path := filepath.Join(dir, "file.txt")

	if err := w.Write(path, []byte("old")); err != nil {
		t.Fatalf("first write failed: %v", err)
	}
	if err := w.Write(path, []byte("new")); err != nil {
		t.Fatalf("second write failed: %v", err)
	}

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile failed: %v", err)
	}
	if string(got) != "new" {
		t.Fatalf("got %q, want %q", got, "new")
	}
}
