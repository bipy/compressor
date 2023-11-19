package handler

import (
	"compressor/constant"
	"compressor/loader"
	"github.com/charmbracelet/log"
)

var CompressHandlerMap = map[string]CompressHandler{}

type CompressHandler interface {
	Name() string
	Compress(raw []byte) ([]byte, bool)
	Parse(cfg *loader.Config)
}

func Register(h CompressHandler) {
	CompressHandlerMap[h.Name()] = h
	log.Debugf("[CompressHandler] registered: %s", h.Name())
}

func GetCompressHandler(name string) CompressHandler {
	if name == "jpg" {
		name = "jpeg"
	}
	h, ok := CompressHandlerMap[name]
	if !ok {
		log.Errorf("[CompressHandler] handler not found, use default. name=%s", name)
		return GetCompressHandler(constant.CompressHandlerJPEG)
	}
	log.Infof("[CompressHandler] use compress handler: %s", name)
	return h
}
