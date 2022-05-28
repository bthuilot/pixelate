package logging

import (
	"io"
	"log"
	"os"
)

const (
	ErrorLevel = Level(iota)
	WarningLevel
	InfoLevel
)

var (
	WarningLogger *log.Logger
	InfoLogger    *log.Logger
	ErrorLogger   *log.Logger
)

type Level uint8

func Init() (err error) {
	logLevel := ErrorLevel
	switch os.Getenv("LOG_LEVEL") {
	case "WARNING":
		logLevel = WarningLevel
	case "INFO":
		logLevel = InfoLevel
	}
	var (
		infoOutput io.Writer
		warnOutput io.Writer
		errOutput  io.Writer
		output     *os.File
	)
	if logFile := os.Getenv("LOG_FILE"); logFile != "" {
		output, err = os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			return err
		}
	} else {
		output = os.Stdout
	}
	infoOutput, warnOutput, errOutput = createWriters(output, logLevel)
	InfoLogger = log.New(infoOutput, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	WarningLogger = log.New(warnOutput, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(errOutput, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	return err
}

func createWriters(file *os.File, level Level) (info, warning, err *os.File) {
	err = file
	devNull, _ := os.Open(os.DevNull)
	info = devNull
	warning = devNull
	if level >= WarningLevel {
		warning = file
	}
	if level >= InfoLevel {
		info = file
	}
	return
}
