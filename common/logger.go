package common

import (
	"log"
	"os"
)

func GetFileLogger(id string) *log.Logger {
	// create log file and init logger
	logFile, err := os.OpenFile(id+".log", os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		panic("Cannot Create Log File")
	}
	return log.New(logFile, "", log.LstdFlags)
}

func GetLogger() *log.Logger {
	return log.New(os.Stdout, "", log.LstdFlags)
}
