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

	flag.IntVar(&cfg.ThreadCount, "j", 4, "thread count")
	flag.StringVar(&cfg.InputPath, "i", "", "input path")
	flag.StringVar(&cfg.OutputPath, "o", "", "output path")
	flag.StringVar(&cfg.OutputType, "t", "jpeg", "output type: jpg/jpeg/png")

	flag.IntVar(&cfg.Quality, "q", 90, "output quality: 0-100")
	flag.IntVar(&cfg.MaxWidth, "width", math.MaxInt, "max image width, default is unlimited")
	flag.IntVar(&cfg.MaxHeight, "height", math.MaxInt, "max image height, default is unlimited")
	flag.StringVar(&cfg.AcceptedInputFormat, "accept", "jpg jpeg png", "accepted input format, default: [jpg jpeg png]")

	flag.Parse()

	if ok := cfg.Check(); !ok {
		os.Exit(1)
	}

	pipeline.Run(cfg)
}
