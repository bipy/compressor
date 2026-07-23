package codec

import (
	"bytes"
	"image"
	"image/color"
	"testing"
)

func createTestImage() *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, 50, 50))
	for y := 0; y < 50; y++ {
		for x := 0; x < 50; x++ {
			img.Set(x, y, color.RGBA{R: uint8(x * 5), G: uint8(y * 5), B: 128, A: 255})
		}
	}
	return img
}

func TestFactory_Encode_JPEG(t *testing.T) {
	f := NewFactory()
	img := createTestImage()

	var buf bytes.Buffer
	err := f.Encode(&buf, img, "jpeg", EncodeOptions{Quality: 80})
	if err != nil {
		t.Fatalf("Encode failed: %v", err)
	}

	if buf.Len() == 0 {
		t.Fatal("encoded data should not be empty")
	}
}

func TestFactory_Encode_PNG(t *testing.T) {
	f := NewFactory()
	img := createTestImage()

	var buf bytes.Buffer
	err := f.Encode(&buf, img, "png", EncodeOptions{})
	if err != nil {
		t.Fatalf("Encode failed: %v", err)
	}

	if buf.Len() == 0 {
		t.Fatal("encoded data should not be empty")
	}
}

func TestFactory_EncodeToBytes(t *testing.T) {
	f := NewFactory()
	img := createTestImage()

	data, err := f.EncodeToBytes(img, "jpeg", EncodeOptions{Quality: 90})
	if err != nil {
		t.Fatalf("EncodeToBytes failed: %v", err)
	}

	if len(data) == 0 {
		t.Fatal("encoded data should not be empty")
	}
}

func TestFactory_Encode_UnsupportedFormat(t *testing.T) {
	f := NewFactory()
	img := createTestImage()

	var buf bytes.Buffer
	err := f.Encode(&buf, img, "gif", EncodeOptions{})
	if err == nil {
		t.Fatal("expected error for unsupported format")
	}
}

func TestFactory_SupportsOutput(t *testing.T) {
	f := NewFactory()

	tests := []struct {
		format string
		want   bool
	}{
		{"jpeg", true},
		{"jpg", true},
		{"png", true},
		{"webp", true},
		{"gif", false},
		{"bmp", false},
		{"tiff", false},
		{"heic", false},
	}

	for _, tt := range tests {
		got := f.SupportsOutput(tt.format)
		if got != tt.want {
			t.Errorf("SupportsOutput(%q) = %v, want %v", tt.format, got, tt.want)
		}
	}
}
