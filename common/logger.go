package common

import (
	"io"
	"log"
	"os"
)

func GetLogger(id string) *log.Logger {
	// create log file and init logger
	logFile, err := os.OpenFile(id+".log", os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		logger := log.New(os.Stdout, "", log.LstdFlags)
		logger.Println(Red("Cannot Create Log File"))
		return logger
	}
	logger := log.New(io.MultiWriter(os.Stdout, logFile), "", log.LstdFlags)
	return logger
}
