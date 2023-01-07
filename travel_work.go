package main

import (
	"compressor/platform"
	"github.com/fatih/color"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func travel() {
	if config.SingleFileMode {
		taskList = append(taskList, Task{
			Input:  config.InputPath,
			Output: strings.TrimSuffix(config.InputPath, filepath.Ext(config.InputPath)) + platform.OutputFormat,
		})
		total = 1
		return
	}

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
			if ext := strings.ToLower(filepath.Ext(d.Name()))[1:]; config.IsAccept(ext) {
				newPath := filepath.Join(config.OutputPath, strings.TrimPrefix(path, config.InputPath))
				newPath = strings.TrimSuffix(newPath, filepath.Ext(newPath)) + platform.OutputFormat
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
