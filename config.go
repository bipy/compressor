package main

import (
	"compressor/common"
	"encoding/json"
	"fmt"
	"image/jpeg"
	"io/ioutil"
	"os"
	"path/filepath"
)

const (
	JpegQualityMax = 100
	JpegQualityMin = 1
	OutputFormat   = ".jpg"
	MaxRenameRetry = 7
)

type Config struct {
	ThreadCount  int             `json:"thread_count"`
	InputFormat  []string        `json:"input_format"`
	InputPath    string          `json:"input_path"`
	OutputPath   string          `json:"output_path"`
	Quality      int             `json:"quality"`
	jpegQuality  *jpeg.Options   // jpeg quality
	acceptFormat map[string]bool // input file format
}

func LoadConfig(configPathPtr *string) *Config {
	// load json
	configFile, err := ioutil.ReadFile(*configPathPtr)
	if err != nil {
		logger.Println(common.Red("Config File Load Failed"))
		os.Exit(1)
	}

	config := &Config{}

	// parse json
	if err := json.Unmarshal(configFile, &config); err != nil {
		logger.Println(common.Red("Json Unmarshal Failed"))
		os.Exit(1)
	}

	ParseConfig(config)

	return config
}

func ParseConfig(config *Config) {
	// check quality
	if config.Quality < JpegQualityMin || config.Quality > JpegQualityMax {
		logger.Println(common.Red("Quality Value Should Between 1 and 100"))
		os.Exit(1)
	}

	// check input path
	config.InputPath = filepath.Clean(config.InputPath)
	if inputInfo, err := os.Stat(config.InputPath); err != nil {
		if os.IsNotExist(err) {
			logger.Println(common.Red("Input Path Not Found"))
			os.Exit(1)
		}
		if !inputInfo.IsDir() {
			logger.Println(common.Red("Input Path Should Be a Directory"))
			os.Exit(1)
		}
	}

	// check output path
	if config.OutputPath != "" {
		if config.InputPath == config.OutputPath {
			logger.Println(common.Red("Output Path Cannot Be Same As Input Path"))
			os.Exit(1)
		}
	} else {
		config.OutputPath = config.InputPath + "_" + id
	}

	config.jpegQuality = &jpeg.Options{Quality: config.Quality}

	// initialize accept input format
	config.acceptFormat = make(map[string]bool)
	for _, v := range config.InputFormat {
		config.acceptFormat[fmt.Sprintf(".%s", v)] = true
	}
}
