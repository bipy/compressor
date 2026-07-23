package transform

import (
	"errors"
	"image"

	"github.com/disintegration/imaging"
)

type Options struct {
	MaxWidth        int
	MaxHeight       int
	AutoOrientation bool
	Orientation     int
}

func Apply(img image.Image, opt Options) (image.Image, error) {
	if img == nil {
		return nil, errors.New("transform: nil image")
	}

	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()
	if width <= 0 || height <= 0 {
		return nil, errors.New("transform: invalid image bounds")
	}

	if opt.AutoOrientation {
		img = applyOrientation(img, opt.Orientation)
		bounds = img.Bounds()
		width, height = bounds.Dx(), bounds.Dy()
	}

	if !needsResize(width, height, opt) {
		return img, nil
	}

	targetWidth, targetHeight := resizeTarget(width, height, opt)
	if targetWidth > 0 {
		return imaging.Resize(img, targetWidth, 0, imaging.Lanczos), nil
	}
	return imaging.Resize(img, 0, targetHeight, imaging.Lanczos), nil
}

func applyOrientation(img image.Image, orientation int) image.Image {
	switch orientation {
	case 2:
		return imaging.FlipH(img)
	case 3:
		return imaging.Rotate180(img)
	case 4:
		return imaging.FlipV(img)
	case 5:
		return imaging.Transpose(img)
	case 6:
		return imaging.Rotate270(img)
	case 7:
		return imaging.Transverse(img)
	case 8:
		return imaging.Rotate90(img)
	default:
		return img
	}
}

func needsResize(width, height int, opt Options) bool {
	if opt.MaxWidth > 0 && width > opt.MaxWidth {
		return true
	}
	if opt.MaxHeight > 0 && height > opt.MaxHeight {
		return true
	}
	return false
}

func resizeTarget(width, height int, opt Options) (int, int) {
	if opt.MaxWidth <= 0 {
		return 0, opt.MaxHeight
	}
	if opt.MaxHeight <= 0 {
		return opt.MaxWidth, 0
	}

	widthRatio := float64(opt.MaxWidth) / float64(width)
	heightRatio := float64(opt.MaxHeight) / float64(height)
	if widthRatio <= heightRatio {
		return opt.MaxWidth, 0
	}
	return 0, opt.MaxHeight
}
