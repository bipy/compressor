package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type Config struct {
	ImageflowTool string `json:"imageflow_tool"`
	ThreadCount   int    `json:"thread_count"`
	InputPath     string `json:"input_path"`
	OutputPath    string `json:"output_path"`
	Command       string
	ProcessType   string
	OutputImage   struct {
		Quality      int    `json:"quality"`
		OutputFormat string `json:"output_format"`
		Resize       struct {
			Enable   bool `json:"enable"`
			ResizeBy int  `json:"resize_by"`
			Width    int  `json:"width"`
			Height   int  `json:"height"`
		} `json:"resize"`
	} `json:"output_image"`
}

type node struct {
	Input  string
	Output string
}

var (
	ID             string // use unix timestamp as process id
	logger         *log.Logger
	config         = Config{} // from json
	total, count   int32      // the number of images
	failList       []node     // gather all failed jobs for summary
	nodeCh, failCh chan node
	wg             = sync.WaitGroup{}
)

func init() {
	// parse args
	flag.Usage = usage
	var configPath string
	flag.StringVar(&configPath, "c", "config.json", "specific configuration")
	flag.Parse()

	if s, err := os.Stat(configPath); err != nil {
		if os.IsNotExist(err) || s.IsDir() {
			fmt.Println("ERROR: Wrong Argument")
			fmt.Println("Use '-h' for help")
			os.Exit(1)
		}
	}

	// init process ID
	ID = strconv.FormatInt(time.Now().Unix(), 10)

	// create log file and init logger
	logFile, err := os.OpenFile(ID+".log", os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		logger = log.New(os.Stdout, "", log.LstdFlags)
		logger.Println("CANNOT CREATE LOG FILE")
	} else {
		logger = log.New(io.MultiWriter(os.Stdout, logFile), "", log.LstdFlags)
	}

	// load json
	configFile, err := ioutil.ReadFile(configPath)
	if err != nil {
		logger.Panicln("LOAD CONFIG FILE FAILED")
	}

	// parse json
	if err := json.Unmarshal(configFile, &config); err != nil {
		logger.Panicln("JSON LOAD FAILED")
	}

	// check tool
	_, err = os.Stat(config.ImageflowTool)
	if err != nil {
		if os.IsNotExist(err) {
			logger.Panicln("imageflow_tool IS NOT FOUND")
		}
	}

	// check input path
	if inputInfo, err := os.Stat(config.InputPath); err != nil {
		if os.IsNotExist(err) {
			logger.Panicln("INPUT PATH IS NOT FOUND")
		}
		if !inputInfo.IsDir() {
			logger.Panicln("INPUT PATH SHOULD BE A DIRECTORY")
		}
	}

	// check output path
	if config.OutputPath != "" {
		if config.InputPath == config.OutputPath {
			logger.Panicln("OUTPUT PATH IS AS SAME AS INPUT PATH")
		}
		if _, err := os.Stat(config.OutputPath); err != nil {
			if os.IsNotExist(err) {
				if e := os.MkdirAll(config.OutputPath, 0755); e != nil {
					logger.Panicln("OUTPUT PATH IS INCORRECT")
				}
			}
		}
	} else {
		config.OutputPath = config.InputPath + "_" + ID
		if e := os.Mkdir(config.OutputPath, 0755); e != nil {
			logger.Panicln("OUTPUT PATH AUTO GENERATE FAILED")
		}
	}

	// initialize process type
	config.ProcessType = "v1/querystring"

	// initialize command
	config.Command = fmt.Sprintf("format=%s&quality=%d", config.OutputImage.OutputFormat, config.OutputImage.Quality)
	if config.OutputImage.Resize.Enable {
		if config.OutputImage.Resize.ResizeBy == 0 {
			config.Command += fmt.Sprintf("&width=%d", config.OutputImage.Resize.Width)
		} else if config.OutputImage.Resize.ResizeBy == 1 {
			config.Command += fmt.Sprintf("&height=%d", config.OutputImage.Resize.Height)
		}
	}

	// initialize channel
	nodeCh = make(chan node, 1024)
	failCh = make(chan node, 1024)
}

// single goroutine!
func travel() {
	// find all images
	err := filepath.Walk(config.InputPath, func(path string, info os.FileInfo, e error) error {
		if e != nil {
			logger.Printf("WALK ERROR: %s %v", path, e)
			return e
		}
		if !info.IsDir() {
			if ext := strings.ToLower(filepath.Ext(info.Name())); ext == ".jpg" || ext == ".png" {
				newPath := filepath.Join(config.OutputPath, path[len(config.InputPath):])
				newPath = newPath[:len(newPath)-3] + config.OutputImage.OutputFormat
				if err := os.MkdirAll(filepath.Dir(newPath), 0755); err != nil {
					logger.Panicln("CREATE NEW PATH FAILED")
				}
				nodeCh <- node{Input: path, Output: newPath}
				total++
			}
		}
		return nil
	})
	if err != nil {
		logger.Panicln("WALK ERROR")
	}

	// close the channel cuz travel is the only sender
	close(nodeCh)
	wg.Done()
}

// compress job, multi-goroutine
func compress() {
	// get job from channel
	// channel nodeCh will be closed by sender
	for j := range nodeCh {
		// execute command and get response
		r, err := exec.Command(
			config.ImageflowTool,
			config.ProcessType,
			"--in",
			j.Input,
			"--out",
			j.Output,
			"--command",
			config.Command).Output()

		// check if response contains success code "200"
		if err != nil || !strings.Contains(string(r), "\"code\": 200") {
			// if failed, push to fail channel (multi-sender)
			failCh <- j
		} else {
			// interface
			// increment and get (CAS)
			v := atomic.LoadInt32(&count)
			for ; !atomic.CompareAndSwapInt32(&count, v, v+1); {
				v = atomic.LoadInt32(&count)
			}
			logger.Printf("(%d/%d) %s -> %s succeed", v+1, total, j.Input, j.Output)
		}
	}

	// push signal when this goroutine is finished
	// use '#' as signal
	// since '#' is not a legal filepath
	failCh <- node{Input: "#", Output: "#"}
	wg.Done()
}

func process() {
	logger.Println("========= Pending =========")

	// two goroutine:
	// 1. travel filepath
	// 2. transfer data from failCh to failList
	wg.Add(2)

	// travel filepath
	go travel()

	// multi-thread compress
	for i := 0; i < config.ThreadCount; i++ {
		wg.Add(1)
		go compress()
	}

	// transfer
	// cuz failCh sender is about to close and buffer is limited!
	go func() {
		// count finished goroutine by finish-signal '#'
		stopCount := 0
		for i := range failCh {
			if i.Input == "#" {
				stopCount++
			} else {
				failList = append(failList, i)
			}
			if stopCount == config.ThreadCount {
				break
			}
		}
		wg.Done()
	}()

	// block main thread until all goroutine is finished
	wg.Wait()
	logger.Println("========= Done =========")
}

func summary() {
	var failCount = int32(len(failList))
	if failCount > 0 {
		logger.Println("Oops! Some of them are failed...")
		for _, n := range failList {
			logger.Printf("Fail: %s", n.Input)
		}
	}
	logger.Println("Process Complete!")
	logger.Printf("Total: %d - Success: %d - Fail: %d",
		total, total-failCount, failCount)
}

func main() {
	process()
	summary()
}

func usage() {
	_, _ = fmt.Fprintf(os.Stderr,
		`Version: 0.1
Usage: compressor [-h] [-c filename]

Options:
  -h
    	show this help
`)
	flag.PrintDefaults()
}
