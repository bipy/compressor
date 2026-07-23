package codec

import (
	"strings"
)

type Factory struct {
	inputs map[string]struct{}
}

func NewFactory() *Factory {
	inputs := map[string]struct{}{
		"jpeg": {},
		"png":  {},
		"gif":  {},
		"bmp":  {},
		"tiff": {},
		"webp": {},
		"heic": {},
	}

	return &Factory{inputs: inputs}
}

func (f *Factory) SupportsInput(format string) bool {
	if f == nil {
		return false
	}
	canonical := normalizeFormat(format)
	if canonical == "" {
		return false
	}
	_, ok := f.inputs[canonical]
	return ok
}

func normalizeFormat(format string) string {
	trimmed := strings.TrimSpace(strings.TrimPrefix(strings.ToLower(format), "."))
	switch trimmed {
	case "jpg":
		return "jpeg"
	case "jpeg", "png", "gif", "bmp", "tif", "tiff", "webp", "heic", "heif":
		if trimmed == "tif" {
			return "tiff"
		}
		if trimmed == "heif" {
			return "heic"
		}
		return trimmed
	default:
		return trimmed
	}
}
