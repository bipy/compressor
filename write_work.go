package main

import (
	"fmt"
	"github.com/fatih/color"
	"os"
)

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
