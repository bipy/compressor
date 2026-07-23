package executor

import (
	"context"
	"image"
	"image/color"
	"image/jpeg"
	"os"
	"path/filepath"
	"testing"

	"compressor/internal/model"
)

func createBenchImage(b *testing.B, path string) {
	b.Helper()
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		b.Fatal(err)
	}
	f, err := os.Create(path)
	if err != nil {
		b.Fatal(err)
	}
	defer f.Close()
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))
	for y := 0; y < 100; y++ {
		for x := 0; x < 100; x++ {
			img.Set(x, y, color.RGBA{R: uint8(x), G: uint8(y), B: 128, A: 255})
		}
	}
	jpeg.Encode(f, img, nil)
}

func BenchmarkExecute_SingleJob(b *testing.B) {
	dir := b.TempDir()
	src := filepath.Join(dir, "input.jpg")
	dst := filepath.Join(dir, "output.jpg")
	createBenchImage(b, src)

	exec := New(Config{Workers: 1, Quality: 80})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
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
		for range results {
		}
	}
}

func BenchmarkExecute_Parallel(b *testing.B) {
	dir := b.TempDir()
	exec := New(Config{Workers: 4, Quality: 80})

	jobs := make(chan model.Job, 10)
	for i := 0; i < 10; i++ {
		src := filepath.Join(dir, "input"+string(rune('a'+i))+".jpg")
		dst := filepath.Join(dir, "output"+string(rune('a'+i))+".jpg")
		createBenchImage(b, src)
		jobs <- model.Job{
			SourcePath: src,
			TargetPath: dst,
			SourceFmt:  "jpeg",
			TargetFmt:  "jpeg",
		}
	}
	close(jobs)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ctx := context.Background()
		results := exec.Execute(ctx, jobs)
		for range results {
		}
	}
}
