package main

import (
	"bufio"
	"bytes"
	"compressor/common"
	"flag"
	"fmt"
	"github.com/fatih/color"
	"image"
	"image/jpeg"
	_ "image/png"
	"io/fs"
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
	id         string      // use unix timestamp as id
	logger     *log.Logger // global logger with color
	fileLogger *log.Logger // global logger without color
	config     *Config     // from json
	failList   []Task      // gather all failed jobs for summary
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
	logToFilePtr := flag.Bool("log", false, "Save Log as File")
	flag.Parse()

	// initialize process id
	id = strconv.FormatInt(time.Now().Unix(), 10)

	if *configPathPtr != "" {
		// parse config file
		config = LoadConfig(*configPathPtr)
	} else {
		config = &Config{
			ThreadCount: *threadCountPtr,
			InputFormat: strings.Split(*inputFormatPtr, " "),
			InputPath:   *inputPathPtr,
			OutputPath:  *outputPathPtr,
			Quality:     *qualityPtr,
			LogToFile:   *logToFilePtr,
		}
	}
	// initialize logger
	logger = common.GetLogger()

	ParseConfig(config)

	dirMutex = &sync.Mutex{}
	wg = &sync.WaitGroup{}

	// initialize channel
	inCh = make(chan Task, config.ThreadCount<<1)
	outCh = make(chan Task, config.ThreadCount<<1)
	failCh = make(chan Task, 1)
}

func travel() {
	// find all images
	err := filepath.WalkDir(config.InputPath, func(path string, d fs.DirEntry, e error) error {
		if e != nil {
			logger.Println(color.RedString("Walk Error:"), path, e.Error())
			if config.LogToFile {
				fileLogger.Println("Walk Error:", path, e.Error())
			}
			return e
		}
		if !d.IsDir() {
			if ext := strings.ToLower(filepath.Ext(d.Name()))[1:]; config.isAccept(ext) {
				newPath := filepath.Join(config.OutputPath, filepath.Base(path))
				newPath = strings.TrimSuffix(newPath, filepath.Ext(newPath)) + OutputFormat
				if err := os.MkdirAll(filepath.Dir(newPath), 0755); err != nil {
					logger.Println(color.RedString("Create New Path Failed"))
					if config.LogToFile {
						fileLogger.Println("Create New Path Failed")
					}
					return err
				}
				taskList = append(taskList, Task{Input: path, Output: newPath})
			}
		}
		return nil
	})
	if err != nil {
		panic(err.Error())
	}
	total = len(taskList)
}

func doTask(t *Task) error {
	file, err := os.Open(t.Input)
	if err != nil {
		return err
	}

	filename, err := common.Touch(t.Output, dirMutex, MaxRenameRetry)
	if err != nil {
		return err
	}
	t.Output = filename

	img, _, err := image.Decode(file)
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	err = jpeg.Encode(buf, img, config.jpegQuality)
	if err != nil {
		return err
	}
	t.Data = buf.Bytes()
	return nil
}

// compress job, multiple goroutine
func compress() {
	defer wg.Done()
	// get job from channel,
	// channel inCh will be closed by sender
	for t := range inCh {
		// check if success
		if err := doTask(&t); err != nil {
			// if failed, push to fail channel (multi-sender)
			t.Data = []byte(err.Error())
			failCh <- t
			continue
		}
		outCh <- t
	}
}

func writeToFiles() {
	defer wg.Done()
	count := 0
	for t := range outCh {
		err := os.WriteFile(t.Output, t.Data, 0644)
		if err != nil {
			t.Data = []byte(err.Error())
			failCh <- t
			continue
		}
		count++
		logger.Println(color.GreenString(fmt.Sprintf("(%d/%d)", count, total)),
			t.Input, color.GreenString("->"), t.Output)
		if config.LogToFile {
			fileLogger.Println(fmt.Sprintf("(%d/%d)", count, total),
				t.Input, "->", t.Output)
		}
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
	for t := range taskList {
		inCh <- taskList[t]
	}
}

func process() {
	// confirm tasks
	logger.Println(color.GreenString("Input Path:"), config.InputPath)
	logger.Println(color.GreenString("Output Path:"), config.OutputPath)
	logger.Println(color.GreenString("Thread Count:"), strconv.Itoa(config.ThreadCount))
	logger.Println(color.GreenString("Accept Format:"), strings.Join(config.InputFormat, ", "))
	logger.Println(color.GreenString("JPEG Quality:"), strconv.Itoa(config.Quality))
	if config.LogToFile {
		logger.Println(color.GreenString("Log:"), id+".log")
	} else {
		logger.Println(color.GreenString("Log:"), "stdout")
	}
	logger.Println(color.YellowString("Continue? (Y/n)"))

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		input := strings.ToLower(scanner.Text())
		if input == "n" {
			logger.Println(color.RedString("Abort!"))
			os.Exit(1)
		} else if input != "" && input != "y" {
			logger.Println(color.YellowString("Continue? (Y/n)"))
		} else {
			break
		}
	}

	if config.LogToFile {
		fileLogger = common.GetFileLogger(id)
		fileLogger.Println("Input Path:", config.InputPath)
		fileLogger.Println("Output Path:", config.OutputPath)
		fileLogger.Println("Thread Count:", strconv.Itoa(config.ThreadCount))
		fileLogger.Println("Accept Format:", strings.Join(config.InputFormat, ", "))
		fileLogger.Println("JPEG Quality:", strconv.Itoa(config.Quality))
		fileLogger.Println("Log:", id+".log")
	}

	// travel filepath
	travel()
	logger.Println(color.GreenString("Found:"), strconv.Itoa(len(taskList)))
	if config.LogToFile {
		fileLogger.Println("Found:", strconv.Itoa(len(taskList)))
	}

	logger.Println(color.BlueString("========= Pending ========="))
	if config.LogToFile {
		fileLogger.Println("========= Pending =========")
	}

	// close by failCh
	go transferFailList()

	// close by outCh
	go writeToFiles()

	// close inCh
	go transferTaskList()

	// multi-thread compress
	for i := 0; i < config.ThreadCount; i++ {
		wg.Add(1)
		go compress()
	}
	// block main thread until all goroutine is finished
	wg.Wait()

	// compress finished
	// close writeToFiles()
	wg.Add(1)
	close(outCh)
	wg.Wait()

	// close transferFailList()
	wg.Add(1)
	close(failCh)
	wg.Wait()

	logger.Println(color.BlueString("=========  Done!  ========="))
	if config.LogToFile {
		fileLogger.Println("=========  Done!  =========")
	}
}

func summary() {
	var failCount = len(failList)
	if failCount > 0 {
		logger.Println(color.YellowString("Oops! Some of them are failed..."))
		if config.LogToFile {
			fileLogger.Println("Oops! Some of them are failed...")
		}
		for _, n := range failList {
			logger.Println(color.RedString("Failed:"), n.Input, "-", string(n.Data))
			if config.LogToFile {
				fileLogger.Println("Failed:", n.Input, "-", string(n.Data))
			}
		}
	}
	logger.Println(color.GreenString("Total:"), total, "-", color.RedString("Failed:"), failCount)
	if config.LogToFile {
		fileLogger.Println("Total:", total, "-", "Failed:", failCount)
	}
}

func main() {
	process()
	summary()
}

func usage() {
	_, _ = fmt.Fprintf(os.Stderr,
		`Version: 2.3
Usage: compressor [-h] [Options]

Options:
  -h
    	show this help
`)
	flag.PrintDefaults()
}
