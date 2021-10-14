package main

import (
    "bytes"
    "compressor/common"
    "encoding/json"
    "flag"
    "fmt"
    "image"
    "image/jpeg"
    _ "image/png"
    "io"
    "io/ioutil"
    "log"
    "os"
    "path/filepath"
    "strconv"
    "strings"
    "sync"
    "sync/atomic"
    "time"
)

type Config struct {
    ThreadCount int      `json:"thread_count"`
    InputFormat []string `json:"input_format"`
    InputPath   string   `json:"input_path"`
    OutputPath  string   `json:"output_path"`
    Quality     int      `json:"quality"`
}

type node struct {
    Input  string
    Output string
}

const OutputFormat = ".jpg"

var (
    id           string // use unix timestamp as id
    logger       *log.Logger
    config       Config // from json
    total, count int32  // the number of images
    acceptFormat map[string]bool
    jpegQuality  *jpeg.Options
    failList     []node // gather all failed jobs for summary
    nodeCh       chan node
    failCh       chan node
    wg           = sync.WaitGroup{}
    dirMutex     = sync.Mutex{}
    travelDone   = false
)

func init() {
    // parse args
    flag.Usage = usage
    configPathPtr := flag.String("c", "config.json", "Configuration Filepath")
    flag.Parse()

    // init process id
    id = strconv.FormatInt(time.Now().Unix(), 10)

    // create log file and init logger
    logFile, err := os.OpenFile(id+".log", os.O_WRONLY|os.O_CREATE, 0644)
    if err != nil {
        logger = log.New(os.Stdout, "", log.LstdFlags)
        logger.Println(common.Red("Cannot Create Log File"))
    } else {
        logger = log.New(io.MultiWriter(os.Stdout, logFile), "", log.LstdFlags)
    }

    // load json
    configFile, err := ioutil.ReadFile(*configPathPtr)
    if err != nil {
        logger.Println(common.Red("Config File Load Failed"))
        os.Exit(1)
    }

    // parse json
    if err := json.Unmarshal(configFile, &config); err != nil {
        logger.Println(common.Red("Json Unmarshal Failed"))
        os.Exit(1)
    }

    // check quality
    if config.Quality <= 0 || config.Quality > 100 {
        logger.Println(common.Red("Quality Value Should Between 1 and 100"))
        os.Exit(1)
    }
    jpegQuality = &jpeg.Options{Quality: config.Quality}

    // check input path
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
        if _, err := os.Stat(config.OutputPath); err != nil {
            if os.IsNotExist(err) {
                if e := os.MkdirAll(config.OutputPath, 0755); e != nil {
                    logger.Println(common.Red("Create Output Path Failed"))
                    os.Exit(1)
                }
            }
        }
    } else {
        config.OutputPath = config.InputPath + "_" + id
        if e := os.Mkdir(config.OutputPath, 0755); e != nil {
            logger.Println(common.Red("Create Output Path Failed"))
            os.Exit(1)
        }
    }

    // initialize accept input format
    acceptFormat = make(map[string]bool)
    for _, v := range config.InputFormat {
        acceptFormat[fmt.Sprintf(".%s", v)] = true
    }

    // initialize channel
    nodeCh = make(chan node, 32768)
    failCh = make(chan node, 32768)
}

func travel() {
    // close the channel cuz travel is the only sender
    defer close(nodeCh)
    defer wg.Done()

    // find all images
    err := filepath.Walk(config.InputPath, func(path string, info os.FileInfo, e error) error {
        if e != nil {
            logger.Printf("%s %s %v", common.Red("Walk Error:"), path, e.Error())
            return e
        }
        if !info.IsDir() {
            if ext := strings.ToLower(filepath.Ext(info.Name())); acceptFormat[ext] {
                newPath := filepath.Join(config.OutputPath, filepath.Base(path))
                newPath = strings.TrimSuffix(newPath, filepath.Ext(newPath)) + OutputFormat
                if err := os.MkdirAll(filepath.Dir(newPath), 0755); err != nil {
                    logger.Println(common.Red("Create New Path Failed"))
                    os.Exit(1)
                }
                nodeCh <- node{Input: path, Output: newPath}
                total++
            }
        }
        return nil
    })
    if err != nil {
        logger.Println(common.Red("Walk Error"))
        os.Exit(1)
    }
    travelDone = true
}

// touch file
func touch(filename *string) bool {
    dirMutex.Lock()
    defer dirMutex.Unlock()
    touched := false
    for i := 0; i < 7; i++ {
        _, err := os.Stat(*filename)
        // if file exist
        if err == nil {
            *filename = fmt.Sprintf("%s (%d)%s", (*filename)[:len(*filename)-len(OutputFormat)], i, OutputFormat)
        } else if os.IsNotExist(err) {
            // touch
            _, err := os.OpenFile(*filename, os.O_RDONLY|os.O_CREATE, 0644)
            if err == nil {
                touched = true
            }
            break
        }
    }
    return touched
}

// compress job, multiple goroutine
func compress() {
    defer wg.Done()
    // get job from channel,
    // channel nodeCh will be closed by sender
    for n := range nodeCh {
        file, err := os.Open(n.Input)
        // check if success
        if err != nil {
            // if failed, push to fail channel (multi-sender)
            failCh <- n
            continue
        }

        if !touch(&n.Output) {
            failCh <- n
            continue
        }

        img, _, err := image.Decode(file)
        if err != nil {
            failCh <- n
            continue
        }

        buf := new(bytes.Buffer)
        err = jpeg.Encode(buf, img, jpegQuality)
        if err != nil {
            failCh <- n
            continue
        }

        err = ioutil.WriteFile(n.Output, buf.Bytes(), 0644)
        if err != nil {
            failCh <- n
            continue
        }

        // interface
        // increment and get (CAS)
        v := atomic.LoadInt32(&count)
        for !atomic.CompareAndSwapInt32(&count, v, v+1) {
            v = atomic.LoadInt32(&count)
        }
        if travelDone {
            logger.Printf("%s %s %s %s",
                common.Purple(fmt.Sprintf("(%d/%d)", v+1, total)),
                n.Input, common.Green("->"), n.Output)
        } else {
            logger.Printf("%s %s %s %s",
                common.Purple(fmt.Sprintf("(%d/loading)", v+1)),
                n.Input, common.Green("->"), n.Output)        }
    }
}

func process() {
    defer close(failCh)

    logger.Println(common.Blue("========= Pending ========="))
    // travel filepath
    wg.Add(1)
    go travel()

    // transfer
    // cuz failCh sender is about to close and buffer is limited
    // when process finished, failCh will be closed
    go func() {
        for n := range failCh {
            failList = append(failList, n)
        }
    }()

    // multi-thread compress
    for i := 0; i < config.ThreadCount; i++ {
        wg.Add(1)
        go compress()
    }

    // block main thread until all goroutine is finished
    wg.Wait()
    logger.Println(common.Blue("========= Done ========="))
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
        `Version: 1.0
Usage: compressor [-h] [-c filename]

Options:
  -h
    	show this help
`)
    flag.PrintDefaults()
}
