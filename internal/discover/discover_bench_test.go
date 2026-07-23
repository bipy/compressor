package discover

import (
	"os"
	"path/filepath"
	"testing"
)

func BenchmarkDiscover(b *testing.B) {
	dir := b.TempDir()

	for i := 0; i < 1000; i++ {
		f, err := os.Create(filepath.Join(dir, "file"+string(rune('a'+i%26))+".jpg"))
		if err != nil {
			b.Fatal(err)
		}
		f.Close()
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for range Discover(dir) {
		}
	}
}

func BenchmarkDiscover_SingleFile(b *testing.B) {
	dir := b.TempDir()
	path := filepath.Join(dir, "test.jpg")
	f, err := os.Create(path)
	if err != nil {
		b.Fatal(err)
	}
	f.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for range Discover(path) {
		}
	}
}
