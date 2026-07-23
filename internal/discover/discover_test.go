package discover

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDiscoverStreamsFilesFromNestedDirectory(t *testing.T) {
	root := t.TempDir()
	targetDir := t.TempDir()

	one := filepath.Join(root, "a", "one.jpg")
	two := filepath.Join(root, "b", "two.png")
	nonRegularTarget := filepath.Join(targetDir, "target.txt")

	mustWriteFile(t, one)
	mustWriteFile(t, two)
	mustWriteFile(t, nonRegularTarget)

	if err := os.Symlink(nonRegularTarget, filepath.Join(root, "b", "link.jpg")); err != nil {
		t.Skipf("symlinks not supported: %v", err)
	}

	var got []string
	for item := range Discover(root) {
		if item.Err != nil {
			t.Fatalf("unexpected error: %v", item.Err)
		}
		got = append(got, item.Path)
	}

	want := []string{one, two}
	if len(got) != len(want) {
		t.Fatalf("expected %d files, got %d", len(want), len(got))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("expected %q at index %d, got %q", want[i], i, got[i])
		}
	}
}

func TestDiscoverStreamsSingleFileInput(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "single.jpg")
	mustWriteFile(t, path)

	var got []string
	for item := range Discover(path) {
		if item.Err != nil {
			t.Fatalf("unexpected error: %v", item.Err)
		}
		got = append(got, item.Path)
	}

	if len(got) != 1 {
		t.Fatalf("expected 1 file, got %d", len(got))
	}
	if got[0] != path {
		t.Fatalf("expected %q, got %q", path, got[0])
	}
}

func TestDiscoverSkipsSymlinkInput(t *testing.T) {
	root := t.TempDir()
	target := filepath.Join(root, "real.jpg")
	link := filepath.Join(root, "link.jpg")
	mustWriteFile(t, target)

	if err := os.Symlink(target, link); err != nil {
		t.Skipf("symlinks not supported: %v", err)
	}

	var got []string
	for item := range Discover(link) {
		if item.Err != nil {
			t.Fatalf("unexpected error: %v", item.Err)
		}
		got = append(got, item.Path)
	}

	if len(got) != 0 {
		t.Fatalf("expected 0 files, got %d", len(got))
	}
}

func TestDiscoverReturnsErrorForMissingRoot(t *testing.T) {
	root := filepath.Join(t.TempDir(), "missing")

	var errCount int
	for item := range Discover(root) {
		if item.Err == nil {
			t.Fatal("expected an error item for missing root")
		}
		errCount++
	}

	if errCount != 1 {
		t.Fatalf("expected 1 error item, got %d", errCount)
	}
}

func mustWriteFile(t *testing.T, path string) {
	t.Helper()

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("create parent directories: %v", err)
	}
	if err := os.WriteFile(path, []byte("test"), 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}
}
