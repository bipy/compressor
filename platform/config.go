package platform

import (
	"encoding/json"
	"image/jpeg"
	"os"
	"path/filepath"
)

const (
	JpegQualityMax = 100
	JpegQualityMin = 1
	OutputFormat   = ".jpg"
)

type Config struct {
	Id             string            `json:"-"` // use unix timestamp as id
	ThreadCount    int               `json:"thread_count"`
	InputFormat    []string          `json:"input_format"`
	InputPath      string            `json:"input_path"`
	OutputPath     string            `json:"output_path"`
	Quality        int               `json:"quality"`
	LogToFile      bool              `json:"log_to_file"`
	JpegQuality    *jpeg.Options     // jpeg quality
	IsAccept       func(string) bool // input file format
	SingleFileMode bool              // single file mode
}

func LoadConfig(configPath string) *Config {
	// load json
	configFile, err := os.ReadFile(configPath)
	if err != nil {
		panic("Config File Load Failed")
	}

	config := &Config{}

	// parse json
	if err := json.Unmarshal(configFile, &config); err != nil {
		panic("Json Unmarshal Failed")
	}

	return config
}

func ParseConfig(config *Config) {
	// check quality
	if config.Quality < JpegQualityMin || config.Quality > JpegQualityMax {
		panic("Quality Value Should Between 1 and 100")
	}

	config.JpegQuality = &jpeg.Options{Quality: config.Quality}

	// check input path & output path
	config.InputPath = filepath.Clean(config.InputPath)
	info, err := os.Stat(config.InputPath)
	if err != nil {
		panic(err.Error())
	}
	if info.IsDir() {
		// dir mode
		// input: dir
		// output: dir_id / out
		config.SingleFileMode = false
		if config.OutputPath != "" {
			if config.InputPath == config.OutputPath {
				panic("Output Path Cannot Be Same As Input Path")
			}
			config.OutputPath = filepath.Clean(config.OutputPath)
		} else {
			config.OutputPath = config.InputPath + "_" + config.Id
		}
	} else {
		// single file mode
		// input: file
		// output: same dir / out
		config.SingleFileMode = true
		if config.OutputPath != "" {
			config.OutputPath = filepath.Clean(config.OutputPath)
		} else {
			config.OutputPath = filepath.Dir(config.InputPath)
		}
	}

	// initialize accept input format
	config.IsAccept = func(s string) (ok bool) {
		for _, v := range config.InputFormat {
			if s == v {
				return true
			}
		}
		return false
	}
}
