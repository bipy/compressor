package handler

import (
	"github.com/charmbracelet/log"
	"github.com/disintegration/imaging"
	"image"
)

type ResizeHandler struct {
	MaxWidth  int
	MaxHeight int
	Filter    imaging.ResampleFilter
}

func (h *ResizeHandler) Resize(img image.Image) (image.Image, bool) {
	width, height := img.Bounds().Dx(), img.Bounds().Dy()
	if width == 0 || height == 0 {
		log.Errorf("[ResizeHandler] image size is 0. width=%d height=%d", width, height)
		return nil, false
	}

	var dst *image.NRGBA
	// if raw image is larger than config
	if width > h.MaxWidth && height > h.MaxHeight {
		// wider
		if float64(h.MaxWidth)/float64(width) < float64(h.MaxHeight)/float64(height) {
			// resize width
			dst = imaging.Resize(img, h.MaxWidth, 0, h.Filter)
		} else { // higher
			// resize height
			dst = imaging.Resize(img, 0, h.MaxHeight, h.Filter)
		}
	} else if width > h.MaxWidth { // wider
		// resize width
		dst = imaging.Resize(img, h.MaxWidth, 0, h.Filter)
	} else if height > h.MaxHeight { // higher
		// resize height
		dst = imaging.Resize(img, 0, h.MaxHeight, h.Filter)
	}

	return dst, true
}
