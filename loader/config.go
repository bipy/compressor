package loader

import (
	"compressor/constant"
	"github.com/charmbracelet/log"
	"math"
	"os"
	"path/filepath"
)

type Config struct {
	ThreadCount         int
	AcceptedInputFormat string
	InputPath           string
	OutputPath          string
	OutputType          string
	Quality             int
	MaxWidth            int
	MaxHeight           int
	SingleFileMode      bool
	//AutoOrientation     bool
}

func (cfg *Config) Check() bool {
	// check quality
	if cfg.Quality < constant.QualityMin || cfg.Quality > constant.QualityMax {
		log.Errorf("[Config] quality value should between %d and %d. quality=%d",
			constant.QualityMin, constant.QualityMax, cfg.Quality)
		return false
	}

	// check shape limit
	if cfg.MaxWidth == 0 {
		cfg.MaxWidth = math.MaxInt
	}
	if cfg.MaxHeight == 0 {
		cfg.MaxHeight = math.MaxInt
	}

	// check input path & output path
	var err error
	cfg.InputPath, err = filepath.Abs(cfg.InputPath)
	if err != nil {
		log.Errorf("[Config] get abs path failed. err=%v", err)
		return false
	}
	info, err := os.Stat(cfg.InputPath)
	if err != nil {
		log.Errorf("[Config] input path is invalid. err=%v", err)
		return false
	}
	if info.IsDir() {
		// dir mode
		// input: dir
		// output: dir_id / out
		cfg.SingleFileMode = false
		if cfg.OutputPath != "" {
			if cfg.InputPath == cfg.OutputPath {
				log.Errorf("[Config] output path == input path.")
				return false
			}
			cfg.OutputPath = filepath.Clean(cfg.OutputPath)
		} else {
			cfg.OutputPath = cfg.InputPath + "-" + constant.ID
		}
	} else {
		// single file mode
		// input: file
		// output: same dir / out
		cfg.SingleFileMode = true
		if cfg.OutputPath != "" {
			cfg.OutputPath = filepath.Clean(cfg.OutputPath)
		} else {
			cfg.OutputPath = filepath.Dir(cfg.InputPath)
		}
	}

	log.Debugf("[Config] check pass. cfg=%+v", cfg)
	return true
}
