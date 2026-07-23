package codec

import (
	"fmt"
	"image"
	"io"

	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/tiff"
	_ "golang.org/x/image/webp"

	_ "github.com/perkeep/heic"
)

func (f *Factory) Decode(r io.Reader) (image.Image, string, error) {
	if r == nil {
		return nil, "", fmt.Errorf("codec: nil reader")
	}

	img, format, err := image.Decode(r)
	if err != nil {
		return nil, "", err
	}
	canonical := normalizeFormat(format)
	if !f.SupportsInput(canonical) {
		return nil, "", fmt.Errorf("codec: unsupported input format %q", format)
	}
	return img, canonical, nil
}

func (f *Factory) DecodeConfig(r io.Reader) (image.Config, string, error) {
	if r == nil {
		return image.Config{}, "", fmt.Errorf("codec: nil reader")
	}

	cfg, format, err := image.DecodeConfig(r)
	if err != nil {
		return image.Config{}, "", err
	}
	canonical := normalizeFormat(format)
	if !f.SupportsInput(canonical) {
		return image.Config{}, "", fmt.Errorf("codec: unsupported input format %q", format)
	}
	return cfg, canonical, nil
}
