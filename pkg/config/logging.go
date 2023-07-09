package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

var logLevels = map[string]logrus.Level{
	"DEBUG": logrus.DebugLevel,
	"INFO":  logrus.InfoLevel,
	"WARN":  logrus.WarnLevel,
	"ERROR": logrus.ErrorLevel,
	"FATAL": logrus.FatalLevel,
}

func InitLogger(level, logFile string, useSTDOUT bool) error {
	if lvl, ok := logLevels[strings.ToUpper(level)]; ok {
		logrus.SetLevel(lvl)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}

	if !useSTDOUT {

		if logFile == "" {
			logFile = "/var/log/pixelate/server.log"
		}

		if file, err := os.Create(logFile); err != nil {
			return fmt.Errorf("unable to open log file")
		} else {
			logrus.SetOutput(file)
		}
	}
	return nil
}
