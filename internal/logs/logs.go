package logs

import (
	"fmt"
	"log"
	"os"

	"github.com/sirupsen/logrus"
	easy "github.com/t-tomalak/logrus-easy-formatter"
)

type Logger struct {
	logController *logrus.Logger
	logFile       *os.File
}

func New(pathToOutputFile string) *Logger {
	const ferr = "internal.logs.New"

	file, err := os.OpenFile(pathToOutputFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, os.ModePerm)
	if err != nil {
		log.Fatalf("%s: cannot use or create logs file: %s", ferr, err.Error())
	}

	return &Logger{
		logController: &logrus.Logger{
			Out:   file,
			Level: logrus.InfoLevel,
			Formatter: &easy.Formatter{
				TimestampFormat: "2006-01-02 15:04:05",
				LogFormat:       "[%lvl%]: %time% - %msg%",
			},
		},
	}
}

func (l *Logger)Close() {
	l.logFile.Close()
}

func (l *Logger)WriteError(err string) {
	l.logController.Error(fmt.Sprintf("%s\n", err))
}

func (l *Logger)WriteInfo(info string) {
	l.logController.Info(fmt.Sprintf("%s\n", info))
}
