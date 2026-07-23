package codec

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/color"
	"image/gif"
	"image/jpeg"
	"image/png"
	"testing"

	"golang.org/x/image/bmp"
	"golang.org/x/image/tiff"
)

func TestFactorySupportsExtendedInputFormats(t *testing.T) {
	f := NewFactory()
	for _, format := range []string{"jpeg", "png", "gif", "bmp", "tiff", "webp", "heic"} {
		if !f.SupportsInput(format) {
			t.Fatalf("expected support for %s", format)
		}
	}
}

func TestFactoryDecodeAndDecodeConfigExtendedFormats(t *testing.T) {
	f := NewFactory()

	fixtures := []struct {
		name   string
		format string
		build  func(*testing.T) []byte
	}{
		{
			name:   "jpeg",
			format: "jpeg",
			build: func(t *testing.T) []byte {
				t.Helper()
				return encodeToBytes(t, func(w *bytes.Buffer) error {
					return jpeg.Encode(w, sampleImage(), nil)
				})
			},
		},
		{
			name:   "png",
			format: "png",
			build: func(t *testing.T) []byte {
				t.Helper()
				return encodeToBytes(t, func(w *bytes.Buffer) error {
					return png.Encode(w, sampleImage())
				})
			},
		},
		{
			name:   "gif",
			format: "gif",
			build: func(t *testing.T) []byte {
				t.Helper()
				return encodeToBytes(t, func(w *bytes.Buffer) error {
					return gif.Encode(w, samplePaletted(), nil)
				})
			},
		},
		{
			name:   "bmp",
			format: "bmp",
			build: func(t *testing.T) []byte {
				t.Helper()
				return encodeToBytes(t, func(w *bytes.Buffer) error {
					return bmp.Encode(w, sampleImage())
				})
			},
		},
		{
			name:   "tiff",
			format: "tiff",
			build: func(t *testing.T) []byte {
				t.Helper()
				return encodeToBytes(t, func(w *bytes.Buffer) error {
					return tiff.Encode(w, sampleImage(), nil)
				})
			},
		},
		{
			name:   "webp",
			format: "webp",
			build: func(t *testing.T) []byte {
				t.Helper()
				return mustDecodeBase64(t, "UklGRiIAAABXRUJQVlA4IBYAAAAwAQCdASoBAAEADsD+JaQAA3AAAAAA")
			},
		},
	}

	for _, fixture := range fixtures {
		t.Run(fixture.name, func(t *testing.T) {
			data := fixture.build(t)

			cfg, format, err := f.DecodeConfig(bytes.NewReader(data))
			if err != nil {
				t.Fatalf("DecodeConfig() error = %v", err)
			}
			if format != fixture.format {
				t.Fatalf("DecodeConfig() format = %q, want %q", format, fixture.format)
			}
			if cfg.Width <= 0 || cfg.Height <= 0 {
				t.Fatalf("DecodeConfig() returned invalid size: %+v", cfg)
			}

			img, format, err := f.Decode(bytes.NewReader(data))
			if err != nil {
				t.Fatalf("Decode() error = %v", err)
			}
			if format != fixture.format {
				t.Fatalf("Decode() format = %q, want %q", format, fixture.format)
			}
			if img.Bounds().Dx() != cfg.Width || img.Bounds().Dy() != cfg.Height {
				t.Fatalf("Decode() bounds = %v, want %dx%d", img.Bounds(), cfg.Width, cfg.Height)
			}
		})
	}
}

func sampleImage() *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, 1, 1))
	img.Set(0, 0, color.RGBA{R: 255, G: 128, B: 64, A: 255})
	return img
}

func samplePaletted() *image.Paletted {
	img := image.NewPaletted(image.Rect(0, 0, 1, 1), color.Palette{color.Black, color.White})
	img.SetColorIndex(0, 0, 1)
	return img
}

func encodeToBytes(t *testing.T, fn func(*bytes.Buffer) error) []byte {
	t.Helper()

	var buf bytes.Buffer
	if err := fn(&buf); err != nil {
		t.Fatalf("encode fixture: %v", err)
	}
	return buf.Bytes()
}

func mustDecodeBase64(t *testing.T, value string) []byte {
	t.Helper()

	data, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		t.Fatalf("decode base64 fixture: %v", err)
	}
	return data
}
