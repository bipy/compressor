package main

import (
	"bufio"
	"compressor/platform"
	"flag"
	"fmt"
	"github.com/fatih/color"
	"github.com/samber/lo"
	"log"
	"math"
	"os"
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
	logger     *log.Logger      // global logger with color
	fileLogger *log.Logger      // global logger without color
	config     *platform.Config // from json
	failList   []Task           // gather all failed jobs for summary
)

// runtime variable
var (
	total      int             // the number of images
	wg         *sync.WaitGroup // thread limit
	taskList   []Task          // task slice
	inCh       chan Task       // in-task channel
	outCh      chan Task       // out-task channel
	failCh     chan Task       // fail channel
	touchMutex *sync.Mutex
)

func init() {
	// parse args
	flag.Usage = usage
	configPathPtr := flag.String("c", "", "Configuration Filepath")
	threadCountPtr := flag.Int("j", 4, "Thread Count")
	inputPathPtr := flag.String("i", "", "Input Path")
	outputPathPtr := flag.String("o", "", "Output Path")
	qualityPtr := flag.Int("q", 90, "JPEG Quality")
	maxWidthPtr := flag.Int("width", math.MaxInt, "Max Image Width")
	maxHeightPtr := flag.Int("height", math.MaxInt, "Max Image Height")
	inputFormatPtr := flag.String("f", "jpg jpeg png", "Input Format")
	logToFilePtr := flag.Bool("log", false, "Save Log as File")
	flag.Parse()

	if *configPathPtr != "" {
		// parse config file
		config = platform.LoadConfig(*configPathPtr)
	} else {
		config = &platform.Config{
			ThreadCount: *threadCountPtr,
			InputFormat: strings.Split(*inputFormatPtr, " "),
			InputPath:   *inputPathPtr,
			OutputPath:  *outputPathPtr,
			Quality:     *qualityPtr,
			MaxWidth:    *maxWidthPtr,
			MaxHeight:   *maxHeightPtr,
			LogToFile:   *logToFilePtr,
		}
	}

	// initialize process id
	config.Id = strconv.FormatInt(time.Now().Unix(), 10)

	platform.ParseConfig(config)

	// initialize logger
	logger = platform.GetLogger()

	wg = &sync.WaitGroup{}

	// initialize channel
	inCh = make(chan Task, config.ThreadCount<<1)
	outCh = make(chan Task, config.ThreadCount<<1)
	failCh = make(chan Task, 1)
	touchMutex = &sync.Mutex{}
}

func process() {
	// confirm tasks
	logger.Println(color.GreenString("Input Path:"), config.InputPath)
	logger.Println(color.GreenString("Output Path:"), config.OutputPath)
	logger.Println(color.GreenString("Thread Count:"), strconv.Itoa(config.ThreadCount))
	logger.Println(color.GreenString("Accept Format:"), strings.Join(config.InputFormat, ", "))
	logger.Println(color.GreenString("JPEG Quality:"), strconv.Itoa(config.Quality))
	logger.Println(color.GreenString("Max Image Width:"),
		lo.Ternary(config.MaxWidth == math.MaxInt, "Unlimited", strconv.Itoa(config.MaxWidth)))
	logger.Println(color.GreenString("Max Image Height:"),
		lo.Ternary(config.MaxHeight == math.MaxInt, "Unlimited", strconv.Itoa(config.MaxHeight)))
	logger.Println(color.GreenString("Log:"),
		lo.Ternary(config.LogToFile, config.Id+".log", "stdout"))

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
		fileLogger = platform.GetFileLogger(config.Id)
		fileLogger.Println("Input Path:", config.InputPath)
		fileLogger.Println("Output Path:", config.OutputPath)
		fileLogger.Println("Thread Count:", strconv.Itoa(config.ThreadCount))
		fileLogger.Println("Accept Format:", strings.Join(config.InputFormat, ", "))
		fileLogger.Println("JPEG Quality:", strconv.Itoa(config.Quality))
		fileLogger.Println("Max Image Width:",
			lo.Ternary(config.MaxWidth == math.MaxInt, "Unlimited", strconv.Itoa(config.MaxWidth)))
		fileLogger.Println("Max Image Height:",
			lo.Ternary(config.MaxHeight == math.MaxInt, "Unlimited", strconv.Itoa(config.MaxHeight)))
		fileLogger.Println("Log:", config.Id+".log")
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
		`Version: 2.8
Usage: compressor [-h] [Options]

Options:
  -h
    	show this help
`)
	flag.PrintDefaults()
}
