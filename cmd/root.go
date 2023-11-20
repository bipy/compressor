package cmd

import (
	"compressor/constant"
	"compressor/loader"
	"compressor/pipeline"
	"flag"
	"fmt"
	"math"
	"os"
)

func Execute() {
	flag.Usage = func() {
		fmt.Printf(constant.FormattedHelpMessage, constant.Version)
		flag.PrintDefaults()
	}

	cfg := &loader.Config{}

	flag.IntVar(&cfg.ThreadCount, "j", 8, "thread count")
	flag.StringVar(&cfg.InputPath, "i", "", "input path")
	flag.StringVar(&cfg.OutputPath, "o", "", "output path")
	flag.StringVar(&cfg.OutputType, "t", "jpg", "output type: jpg/jpeg/png")

	flag.IntVar(&cfg.Quality, "q", 90, "output quality: 0-100")
	flag.IntVar(&cfg.MaxWidth, "width", math.MaxInt, "max image width")
	flag.IntVar(&cfg.MaxHeight, "height", math.MaxInt, "max image height")
	flag.StringVar(&cfg.AcceptedInputFormat, "accept", "jpg jpeg png", "accepted input format")

	flag.Parse()

	if ok := cfg.Check(); !ok {
		os.Exit(1)
	}

	pipeline.Run(cfg)
}
