package codec

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"

	"github.com/nickalie/go-webpbin"
)

type EncodeOptions struct {
	Quality int // 1-100, used for JPEG and WebP
}

func (f *Factory) Encode(w io.Writer, img image.Image, format string, opts EncodeOptions) error {
	if w == nil {
		return fmt.Errorf("codec: nil writer")
	}
	if img == nil {
		return fmt.Errorf("codec: nil image")
	}

	canonical := normalizeFormat(format)
	if !f.SupportsOutput(canonical) {
		return fmt.Errorf("codec: unsupported output format %q", format)
	}

	quality := opts.Quality
	if quality <= 0 {
		quality = 80
	}

	switch canonical {
	case "jpeg":
		return jpeg.Encode(w, img, &jpeg.Options{Quality: quality})
	case "png":
		return png.Encode(w, img)
	case "webp":
		return encodeWebP(w, img, quality)
	default:
		return fmt.Errorf("codec: unsupported output format %q", format)
	}
}

func (f *Factory) EncodeToBytes(img image.Image, format string, opts EncodeOptions) ([]byte, error) {
	var buf bytes.Buffer
	if err := f.Encode(&buf, img, format, opts); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (f *Factory) SupportsOutput(format string) bool {
	if f == nil {
		return false
	}
	canonical := normalizeFormat(format)
	switch canonical {
	case "jpeg", "png", "webp":
		return true
	default:
		return false
	}
}

func encodeWebP(w io.Writer, img image.Image, quality int) error {
	enc := &webpbin.Encoder{Quality: uint(quality)}
	return enc.Encode(w, img)
}
