// Copyright 2021 Nitric Technologies Pty Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package logger

import (
	"fmt"
	"io"
	"log"

	"github.com/nitrictech/nitric/core/pkg/env"
)

type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
	FATAL
)

var levelNames = [...]string{"DEBUG", "INFO", "WARN", "ERROR", "FATAL"}

func (l LogLevel) String() string {
	return levelNames[l]
}

func LogLevelFromString(level string) LogLevel {
	for i, name := range levelNames {
		if name == level {
			return LogLevel(i)
		}
	}
	return INFO // default to INFO
}

var logLevel = LogLevelFromString(env.LOG_LEVEL.String())

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
