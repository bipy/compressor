package main

import (
	"bufio"
	"bytes"
	"compressor/common"
	"flag"
	"fmt"
	"image"
	"image/jpeg"
	_ "image/png"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Task struct {
	Input  string
	Output string
	Data   []byte
}

var (
	id       string      // use unix timestamp as id
	logger   *log.Logger // global logger
	config   *Config     // from json
	failList []Task      // gather all failed jobs for summary
)

// runtime variable
var (
	total    int             // the number of images
	wg       *sync.WaitGroup // thread limit
	dirMutex *sync.Mutex     // dir lock for creating file
	taskList []Task          // task slice
	inCh     chan Task       // in-task channel
	outCh    chan Task       // out-task channel
	failCh   chan Task       // fail channel

)

func init() {
	// parse args
	flag.Usage = usage
	configPathPtr := flag.String("c", "", "Configuration Filepath")
	threadCountPtr := flag.Int("j", 4, "Thread Count")
	inputPathPtr := flag.String("i", "", "Input Path")
	outputPathPtr := flag.String("o", "", "Output Path")
	qualityPtr := flag.Int("q", 90, "JPEG Quality")
	inputFormatPtr := flag.String("f", "jpg jpeg png", "Input Format")
	flag.Parse()

	// initialize process id
	id = strconv.FormatInt(time.Now().Unix(), 10)

	// initialize logger
	logger = common.GetLogger(id)

	if *configPathPtr != "" {
		// parse config file
		config = LoadConfig(configPathPtr)
	} else {
		config = &Config{
			ThreadCount:  *threadCountPtr,
			InputFormat:  strings.Split(*inputFormatPtr, " "),
			InputPath:    *inputPathPtr,
			OutputPath:   *outputPathPtr,
			Quality:      *qualityPtr,
			jpegQuality:  nil,
			acceptFormat: nil,
		}
		ParseConfig(config)
	}

	dirMutex = &sync.Mutex{}
	wg = &sync.WaitGroup{}

	// initialize channel
	inCh = make(chan Task, config.ThreadCount<<1)
	outCh = make(chan Task, config.ThreadCount<<1)
	failCh = make(chan Task, 1)
}

func travel() {
	// find all images
	err := filepath.Walk(config.InputPath, func(path string, info os.FileInfo, e error) error {
		if e != nil {
			logger.Printf("%s %s %v", common.Red("Walk Error:"), path, e.Error())
			return e
		}
		if !info.IsDir() {
			if ext := strings.ToLower(filepath.Ext(info.Name())); config.acceptFormat[ext] {
				newPath := filepath.Join(config.OutputPath, path[len(config.InputPath):])
				newPath = strings.TrimSuffix(newPath, filepath.Ext(newPath)) + OutputFormat
				if err := os.MkdirAll(filepath.Dir(newPath), 0755); err != nil {
					logger.Println(common.Red("Create New Path Failed"))
					os.Exit(1)
				}
				taskList = append(taskList, Task{Input: path, Output: newPath})
			}
		}
		return nil
	})
	if err != nil {
		logger.Println(common.Red("Walk Error"))
		os.Exit(1)
	}
	total = len(taskList)
}

// compress job, multiple goroutine
func compress() {
	defer wg.Done()
	// get job from channel,
	// channel nodeCh will be closed by sender
	for t := range inCh {
		file, err := os.Open(t.Input)
		// check if success
		if err != nil {
			// if failed, push to fail channel (multi-sender)
			failCh <- t
			continue
		}

		if !common.Touch(&t.Output, dirMutex, MaxRenameRetry) {
			failCh <- t
			continue
		}

		img, _, err := image.Decode(file)
		if err != nil {
			failCh <- t
			continue
		}

		buf := new(bytes.Buffer)
		err = jpeg.Encode(buf, img, config.jpegQuality)
		if err != nil {
			failCh <- t
			continue
		}

		t.Data = buf.Bytes()
		outCh <- t
	}
}

func writeToFiles() {
	defer wg.Done()
	count := 0
	for t := range outCh {
		err := ioutil.WriteFile(t.Output, t.Data, 0644)
		if err != nil {
			failCh <- t
			continue
		}
		count++
		logger.Printf("%s %s %s %s",
			common.Green(fmt.Sprintf("(%d/%d)", count, total)),
			t.Input, common.Green("->"), t.Output)
	}
}

// transfer
// when process finished, failCh will be closed
func transferFailList() {
	defer wg.Done()
	for t := range failCh {
		failList = append(failList, t)
	}
}

// transfer
// close the channel cuz this is the only sender
func transferTaskList() {
	defer close(inCh)
	defer wg.Done()

	for t := range taskList {
		inCh <- taskList[t]
	}
}

func process() {
	// confirm tasks
	logger.Println(common.Green("Input Path:"), config.InputPath)
	logger.Println(common.Green("Output Path:"), config.OutputPath)
	logger.Println(common.Green("Thread Count:"), strconv.Itoa(config.ThreadCount))
	logger.Println(common.Green("Accept Format:"), strings.Join(config.InputFormat, " "))
	logger.Println(common.Green("JPEG Quality:"), strconv.Itoa(config.Quality))
	logger.Println(common.Yellow("Continue? (Y/n)"))
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		input := strings.ToLower(scanner.Text())
		if input == "n" {
			logger.Println(common.Red("Abort!"))
			os.Exit(1)
		} else if input != "" && input != "y" {
			logger.Println(common.Yellow("Continue? (Y/n)"))
		} else {
			break
		}
	}
	if _, err := os.Stat(config.OutputPath); err != nil {
		if os.IsNotExist(err) {
			if e := os.MkdirAll(config.OutputPath, 0755); e != nil {
				logger.Println(common.Red("Create Output Path Failed"))
				os.Exit(1)
			}
		}
	}
	// travel filepath
	travel()
	logger.Println(common.Green("Found:"), strconv.Itoa(len(taskList)))

	logger.Println(common.Blue("========= Pending ========="))
	go transferFailList()
	go writeToFiles()

	wg.Add(1)
	go transferTaskList()

	// multi-thread compress
	for i := 0; i < config.ThreadCount; i++ {
		wg.Add(1)
		go compress()
	}
	// block main thread until all goroutine is finished
	wg.Wait()

	// close by order
	// close writeToFiles()
	wg.Add(1)
	close(outCh)
	wg.Wait()

	// close transferFailList()
	wg.Add(1)
	close(failCh)
	wg.Wait()

	logger.Println(common.Blue("=========  Done!  ========="))
}

func summary() {
	var failCount = len(failList)
	if failCount > 0 {
		logger.Println(common.Yellow("Oops! Some of them are failed..."))
		for _, n := range failList {
			logger.Printf("%s %s", common.Red("Failed:"), n.Input)
		}
	}
	logger.Printf("%s %d - %s %d", common.Green("Total:"), total, common.Red("Failed:"), failCount)
}

func main() {
	process()
	summary()
}

func usage() {
	_, _ = fmt.Fprintf(os.Stderr,
		`Version: 2.2
Usage: compressor [-h] [Options]

Options:
  -h
    	show this help
`)
	flag.PrintDefaults()
}
