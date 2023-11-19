package handler

import (
	"bytes"
	"compressor/constant"
	"compressor/loader"
	"github.com/charmbracelet/log"
	"github.com/disintegration/imaging"
	"image/png"
	"math"
)

type PNGHandler struct {
	CompressionLevel png.CompressionLevel
	AutoOrientation  bool
	Resizer          *ResizeHandler
}

func init() {
	Register(&PNGHandler{CompressionLevel: 0, AutoOrientation: true})
}

func (h *PNGHandler) Name() string {
	return constant.CompressHandlerPNG
}

func (h *PNGHandler) Compress(raw []byte) (data []byte, ok bool) {
	img, err := imaging.Decode(bytes.NewReader(raw), imaging.AutoOrientation(h.AutoOrientation))
	if err != nil {
		log.Errorf("[PNGHandler] decode failed. err=%v", err)
		return nil, false
	}

	if h.Resizer != nil {
		img, ok = h.Resizer.Resize(img)
		if !ok {
			log.Errorf("[PNGHandler] resize failed.")
			return nil, false
		}
	}

	buf := new(bytes.Buffer)
	err = imaging.Encode(buf, img, imaging.PNG, imaging.PNGCompressionLevel(h.CompressionLevel))
	if err != nil {
		log.Errorf("[PNGHandler] encode failed. err=%v", err)
		return nil, false
	}
	return buf.Bytes(), true
}

func (h *PNGHandler) Parse(cfg *loader.Config) {
	if cfg.MaxWidth != math.MaxInt || cfg.MaxHeight != math.MaxInt {
		h.Resizer = &ResizeHandler{
			MaxWidth:  cfg.MaxWidth,
			MaxHeight: cfg.MaxHeight,
			Filter:    imaging.Lanczos,
		}
	}
}
