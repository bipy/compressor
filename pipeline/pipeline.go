package pipeline

import (
	"compressor/handler"
	"compressor/loader"
	"compressor/utils"
	"github.com/charmbracelet/log"
	"os"
	"strings"
	"sync"
	"sync/atomic"
)

func Run(cfg *loader.Config) {
	compressHandler := handler.GetCompressHandler(cfg.OutputType)

	fileHandler := handler.FileHandler{
		Mu:                  &sync.Mutex{},
		BasePath:            cfg.InputPath,
		OutputPath:          cfg.OutputPath,
		AcceptedInputFormat: strings.Split(cfg.AcceptedInputFormat, " "),
		OutputFormat:        cfg.OutputType,
	}

	fileList, dstList := fileHandler.Travel(cfg.SingleFileMode)
	err := utils.Clone(dstList)
	if err != nil {
		log.Fatalf("[FileHandler] make directory failed. err=%v", err)
	}
	total := len(fileList)
	log.Infof("Found: %d", total)
	successCount := &atomic.Uint64{}

	compress := func(idx int) {
		in, out := fileList[idx], dstList[idx]
		raw, e := os.ReadFile(in)
		if e != nil {
			log.Errorf("[Worker] read file failed. err=%v in=%s out=%s", e, in, out)
			return
		}
		res, ok := compressHandler.Compress(raw)
		if !ok {
			log.Errorf("[Worker] compress failed. in=%s out=%s", in, out)
			return
		}
		out, ok = fileHandler.Touch(out)
		if !ok {
			log.Errorf("[Worker] touch file failed. in=%s out=%s", in, out)
			return
		}
		ok = fileHandler.Write(out, res)
		if !ok {
			log.Errorf("[Worker] write file failed. in=%s out=%s", in, out)
			return
		}
		cur := successCount.Add(1)
		log.Infof("[%d/%d] (%s) -> (%s)", cur, total, in, out)
	}

	workerHandler := handler.GetWorkerHandler(cfg.ThreadCount)
	for i := 0; i < total; i++ {
		workerHandler.Run(i, compress)
	}
	workerHandler.Wait()

	success := int(successCount.Load())
	log.Infof("[Summary] success: %d - fail: %d", success, total-success)
}
