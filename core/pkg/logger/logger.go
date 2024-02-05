package logger

import (
	"fmt"
	"io"
	"log"
)

type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
	FATAL
)

func (l LogLevel) String() string {
	return [...]string{"DEBUG", "INFO", "WARN", "ERROR", "FATAL"}[l]
}

var logLevel = INFO

func SetLogLevel(level LogLevel) {
	logLevel = level
}

func GetLogLevel() LogLevel {
	return logLevel
}

func SetOutput(w io.Writer) {
	log.SetOutput(w)
}

func logMsg(level LogLevel, msg string) {
	if level >= logLevel {
		if level == FATAL {
			log.Fatalf("%s %s\n", level.String(), msg)
		} else {
			log.Printf("%s %s\n", level.String(), msg)
		}
	}
}

func Debug(msg string) {
	logMsg(DEBUG, msg)
}

func Debugf(format string, v ...interface{}) {
	Debug(fmt.Sprintf(format, v...))
}

func Info(msg string) {
	logMsg(INFO, msg)
}

func Infof(format string, v ...interface{}) {
	Info(fmt.Sprintf(format, v...))
}

func Warn(msg string) {
	logMsg(WARN, msg)
}

func Warnf(format string, v ...interface{}) {
	Warn(fmt.Sprintf(format, v...))
}

func Error(msg string) {
	logMsg(ERROR, msg)
}

func Errorf(format string, v ...interface{}) {
	Error(fmt.Sprintf(format, v...))
}

func Fatal(msg string) {
	logMsg(FATAL, msg)
}

func Fatalf(format string, v ...interface{}) {
	Fatal(fmt.Sprintf(format, v...))
}
