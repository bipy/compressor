package planner

import (
	"image"
	"image/jpeg"
	"os"
	"path/filepath"
	"testing"
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
	jpeg.Encode(f, img, nil)
}

func BenchmarkPlan(b *testing.B) {
	dir := b.TempDir()
	src := filepath.Join(dir, "test.jpg")
	createBenchImage(b, src)

	p := New("jpeg", nil, dir)
	dest := Dir(filepath.Join(dir, "output"))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p.Plan(src, dest)
	}
}

func BenchmarkPlan_WithCollision(b *testing.B) {
	dir := b.TempDir()
	src := filepath.Join(dir, "test.jpg")
	createBenchImage(b, src)

	outDir := filepath.Join(dir, "output")
	os.MkdirAll(outDir, 0755)

	p := New("jpeg", nil, dir)
	dest := Dir(outDir)

	for i := 0; i < 100; i++ {
		f, _ := os.Create(filepath.Join(outDir, "test.jpg"))
		f.Close()
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p.Plan(src, dest)
	}
}
