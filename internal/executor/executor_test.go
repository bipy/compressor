package executor

import (
	"bytes"
	"context"
	"image"
	"image/color"
	"image/jpeg"
	"os"
	"path/filepath"
	"testing"

	"compressor/internal/model"
)

func createTestImage(t *testing.T, path string) {
	t.Helper()
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	img := image.NewRGBA(image.Rect(0, 0, 100, 100))
	for y := 0; y < 100; y++ {
		for x := 0; x < 100; x++ {
			img.Set(x, y, color.RGBA{R: uint8(x), G: uint8(y), B: 128, A: 255})
		}
	}

	f, err := os.Create(path)
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	defer f.Close()

	if err := jpeg.Encode(f, img, nil); err != nil {
		t.Fatalf("encode: %v", err)
	}
}

func TestExecutor_ProcessJob(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "input.jpg")
	dst := filepath.Join(dir, "output.jpg")
	createTestImage(t, src)

	exec := New(Config{
		Workers: 2,
		Quality: 80,
	})

	jobs := make(chan model.Job, 1)
	jobs <- model.Job{
		SourcePath: src,
		TargetPath: dst,
		SourceFmt:  "jpeg",
		TargetFmt:  "jpeg",
	}
	close(jobs)

	ctx := context.Background()
	results := exec.Execute(ctx, jobs)

	var count int
	for r := range results {
		count++
		if r.Status != model.StatusDone {
			t.Fatalf("expected status done, got %s: %v", r.Status, r.Err)
		}
		if r.InputBytes == 0 {
			t.Fatal("expected non-zero input bytes")
		}
		if r.OutputBytes == 0 {
			t.Fatal("expected non-zero output bytes")
		}
	}

	if count != 1 {
		t.Fatalf("expected 1 result, got %d", count)
	}

	if _, err := os.Stat(dst); err != nil {
		t.Fatalf("output file should exist: %v", err)
	}
}

func TestExecutor_MultipleJobs(t *testing.T) {
	dir := t.TempDir()
	exec := New(Config{Workers: 2, Quality: 80})

	jobs := make(chan model.Job, 3)
	for i := 0; i < 3; i++ {
		src := filepath.Join(dir, "input"+string(rune('a'+i))+".jpg")
		dst := filepath.Join(dir, "output"+string(rune('a'+i))+".jpg")
		createTestImage(t, src)
		jobs <- model.Job{
			SourcePath: src,
			TargetPath: dst,
			SourceFmt:  "jpeg",
			TargetFmt:  "jpeg",
		}
	}
	close(jobs)

	ctx := context.Background()
	results := exec.Execute(ctx, jobs)

	var count int
	for r := range results {
		count++
		if r.Status != model.StatusDone {
			t.Fatalf("job %s failed: %v", r.Job.SourcePath, r.Err)
		}
	}

	if count != 3 {
		t.Fatalf("expected 3 results, got %d", count)
	}
}

func TestExecutor_MissingFile(t *testing.T) {
	dir := t.TempDir()
	exec := New(Config{Workers: 1, Quality: 80})

	jobs := make(chan model.Job, 1)
	jobs <- model.Job{
		SourcePath: filepath.Join(dir, "nonexistent.jpg"),
		TargetPath: filepath.Join(dir, "output.jpg"),
		SourceFmt:  "jpeg",
		TargetFmt:  "jpeg",
	}
	close(jobs)

	ctx := context.Background()
	results := exec.Execute(ctx, jobs)

	r := <-results
	if r.Status != model.StatusFailed {
		t.Fatalf("expected failed status, got %s", r.Status)
	}
	if r.Err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestBytesReaderWrapper(t *testing.T) {
	data := []byte("hello")
	r := bytesReader(data)

	buf := make([]byte, 3)
	n, err := r.Read(buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 3 || !bytes.Equal(buf[:n], []byte("hel")) {
		t.Fatalf("got %q, want %q", buf[:n], "hel")
	}

	n, err = r.Read(buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 2 || !bytes.Equal(buf[:n], []byte("lo")) {
		t.Fatalf("got %q, want %q", buf[:n], "lo")
	}
}
