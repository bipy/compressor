package main

import (
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
    "sync/atomic"
    "time"
)

type node struct {
    Input  string
    Output string
}

var (
    id       string      // use unix timestamp as id
    logger   *log.Logger // global logger
    config   *Config     // from json
    failList []node      // gather all failed jobs for summary
)

// runtime variable
var (
    total, count int32           // the number of images
    wg           *sync.WaitGroup // thread limit
    dirMutex     *sync.Mutex     // dir lock for creating file
    travelDone   bool            // if travel is finished
    nodeCh       chan node       // task channel
    failCh       chan node       // task channel
)

func init() {
    // parse args
    flag.Usage = usage
    configPathPtr := flag.String("c", "config.json", "Configuration Filepath")
    flag.Parse()

    // initialize process id
    id = strconv.FormatInt(time.Now().Unix(), 10)

    // initialize logger
    logger = common.GetLogger(id)

    // parse config file
    config = ParseConfig(configPathPtr)

    travelDone = false
    dirMutex = &sync.Mutex{}
    wg = &sync.WaitGroup{}

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
            if ext := strings.ToLower(filepath.Ext(info.Name())); config.acceptFormat[ext] {
                newPath := filepath.Join(config.OutputPath, path[len(config.InputPath):])
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

        if !common.Touch(&n.Output, dirMutex, MaxRenameRetry) {
            failCh <- n
            continue
        }

        img, _, err := image.Decode(file)
        if err != nil {
            failCh <- n
            continue
        }

        buf := new(bytes.Buffer)
        err = jpeg.Encode(buf, img, config.jpegQuality)
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
                n.Input, common.Green("->"), n.Output)
        }
    }
}

func process() {
    defer close(failCh)

    // transfer
    // when process finished, failCh will be closed
    go func() {
        for n := range failCh {
            failList = append(failList, n)
        }
    }()

    logger.Println(common.Blue("========= Pending ========="))
    // travel filepath
    wg.Add(1)
    go travel()

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
        `Version: 2.1
Usage: compressor [-h] [-c filename]

Options:
  -h
    	show this help
`)
    flag.PrintDefaults()
}
