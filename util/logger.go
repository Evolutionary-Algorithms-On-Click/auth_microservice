package util

import (
	"log"
	"os"
	"time"
)

type LogLevel int

const (
	INFO LogLevel = iota
	WARN
	ERROR
)

type Logger struct {
	logger *log.Logger
	level  LogLevel
}

func NewLogger() *Logger {
	return &Logger{
		logger: log.New(os.Stdout, "", 0),
		level:  INFO,
	}
}

func (l *Logger) Log(level LogLevel, msg string) {
	if level >= l.level {
		timestamp := time.Now().Format(time.DateTime)
		levelStr := ""
		switch level {
		case INFO:
			levelStr = "INFO"
		case WARN:
			levelStr = "WARN"
		case ERROR:
			levelStr = "ERROR"
		}
		l.logger.Printf("%s [%s]: %s", timestamp, levelStr, msg)
	}
}

func (l *Logger) Info(msg string) {
	l.Log(INFO, msg)
}

func (l *Logger) Warn(msg string) {
	l.Log(WARN, msg)
}

func (l *Logger) Error(msg string) {
	l.Log(ERROR, msg)
}

var Log_var =  NewLogger()
