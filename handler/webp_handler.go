package handler

import (
	"bytes"
	"compressor/constant"
	"compressor/loader"
	"github.com/charmbracelet/log"
	"github.com/disintegration/imaging"
	"github.com/nickalie/go-webpbin"
	"image"
	"math"
)

type WEBPHandler struct {
	AutoOrientation bool
	Resizer         *ResizeHandler
	encoder         *webpbin.Encoder
}

func init() {
	Register(&WEBPHandler{AutoOrientation: true})
}

func (h *WEBPHandler) Name() string {
	return constant.CompressHandlerWEBP
}

func (h *WEBPHandler) Compress(raw []byte) (data []byte, ok bool) {
	img, err := imaging.Decode(bytes.NewReader(raw), imaging.AutoOrientation(h.AutoOrientation))
	if err != nil {
		log.Errorf("[WEBPHandler] decode failed. err=%v", err)
		return nil, false
	}

	if h.Resizer != nil {
		img, ok = h.Resizer.Resize(img)
		if !ok {
			log.Errorf("[WEBPHandler] resize failed.")
			return nil, false
		}
	}

	buf := new(bytes.Buffer)
	err = h.encoder.Encode(buf, img)
	if err != nil {
		log.Errorf("[WEBPHandler] encode failed. err=%v", err)
		return nil, false
	}
	return buf.Bytes(), true
}

func (h *WEBPHandler) Parse(cfg *loader.Config) {
	if cfg.MaxWidth != math.MaxInt || cfg.MaxHeight != math.MaxInt {
		h.Resizer = &ResizeHandler{
			MaxWidth:  cfg.MaxWidth,
			MaxHeight: cfg.MaxHeight,
			Filter:    imaging.Lanczos,
		}
	}
	h.encoder = &webpbin.Encoder{Quality: uint(cfg.Quality)}
	// warmup
	err := h.encoder.Encode(new(bytes.Buffer), image.NewNRGBA(image.Rect(0, 0, 1, 1)))
	if err != nil {
		log.Fatalf("[WEBPHandler] warmup failed. err=%v", err)
	} else {
		log.Infof("[WEBPHandler] warmup finished. quality=%d", cfg.Quality)
	}
}
