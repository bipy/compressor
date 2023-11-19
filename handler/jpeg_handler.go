package handler

import (
	"bytes"
	"compressor/constant"
	"compressor/loader"
	"github.com/charmbracelet/log"
	"github.com/disintegration/imaging"
	"math"
)

type JPEGHandler struct {
	Quality         int
	AutoOrientation bool
	Resizer         *ResizeHandler
}

func init() {
	Register(&JPEGHandler{Quality: 90, AutoOrientation: true})
}

func (h *JPEGHandler) Name() string {
	return constant.CompressHandlerJPEG
}

func (h *JPEGHandler) Compress(raw []byte) (data []byte, ok bool) {
	img, err := imaging.Decode(bytes.NewReader(raw), imaging.AutoOrientation(h.AutoOrientation))
	if err != nil {
		log.Errorf("[JPEGHandler] decode failed. err=%v", err)
		return nil, false
	}

	if h.Resizer != nil {
		img, ok = h.Resizer.Resize(img)
		if !ok {
			log.Errorf("[JPEGHandler] resize failed.")
			return nil, false
		}
	}

	buf := new(bytes.Buffer)
	err = imaging.Encode(buf, img, imaging.JPEG, imaging.JPEGQuality(h.Quality))
	if err != nil {
		log.Errorf("[JPEGHandler] encode failed. err=%v", err)
		return nil, false
	}
	return buf.Bytes(), true
}

func (h *JPEGHandler) Parse(cfg *loader.Config) {
	if cfg.MaxWidth != math.MaxInt || cfg.MaxHeight != math.MaxInt {
		h.Resizer = &ResizeHandler{
			MaxWidth:  cfg.MaxWidth,
			MaxHeight: cfg.MaxHeight,
			Filter:    imaging.Lanczos,
		}
	}
	h.Quality = cfg.Quality
}
