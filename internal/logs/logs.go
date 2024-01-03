package logs

import (
	"fmt"
	"log"
	"os"

	"github.com/sirupsen/logrus"
	easy "github.com/t-tomalak/logrus-easy-formatter"
)

type Logger struct {
	logsToFile  *logrus.Logger
	logsToStd   *logrus.Logger
	logFile     *os.File
	silenceMode bool
}

func New(pathToOutputFile string) *Logger {
	const ferr = "internal.logs.New"

	file, err := os.OpenFile(pathToOutputFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, os.ModePerm)
	if err != nil {
		log.Fatalf("%s: cannot use or create log file: %s", ferr, err.Error())
	}

	return &Logger{
		logsToFile: &logrus.Logger{
			Out:   file,
			Level: logrus.InfoLevel,
			Formatter: &easy.Formatter{
				TimestampFormat: "2006-01-02 15:04:05",
				LogFormat:       "[%lvl%]: %time% - %msg%",
			},
		},
		logFile: file,
		logsToStd: &logrus.Logger{
			Out:   os.Stdout,
			Level: logrus.InfoLevel,
			Formatter: &easy.Formatter{
				TimestampFormat: "2006-01-02 15:04:05",
				LogFormat:       "[%lvl%]: %time% - %msg%",
			},
		},
		silenceMode: false,
	}
}

func (l *Logger) Close() {
	l.logFile.Close()
}

func (l *Logger) WriteError(ferr, err string) {
	if l.silenceMode == false {
		l.logsToStd.Error(fmt.Sprintf("%s: %s\n", ferr, err))
	}
	l.logsToFile.Error(fmt.Sprintf("%s: %s\n", ferr, err))
}

func (l *Logger) WriteInfo(info string) {
	if l.silenceMode == false {
		l.logsToStd.Info(fmt.Sprintf("%s\n", info))
	}
	l.logsToFile.Info(fmt.Sprintf("%s\n", info))
}

func (l *Logger) SilenceOperatingMode(mode bool) {
	l.silenceMode = mode
}
