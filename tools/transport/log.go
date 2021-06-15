package main

import (
	"fmt"
	"log"
)

const (
	//LogLevelDebug is debug log level
	LogLevelDebug = 40
	//LogLevelInfo is info log level
	LogLevelInfo = 30
	//LogLevelWarn is warn log level
	LogLevelWarn = 20
	//LogLevelError is error log level
	LogLevelError = 10
)

var logLevel = LogLevelInfo

//SetLogLevel is set log level to l
func SetLogLevel(l int) {
	if l > 0 {
		logLevel = l
	}
}

//DebugLog is the debug level log
func DebugLog(format string, args ...interface{}) {
	if logLevel < LogLevelDebug {
		return
	}
	log.Output(2, fmt.Sprintf("D "+format, args...))
}

//InfoLog is the info level log
func InfoLog(format string, args ...interface{}) {
	if logLevel < LogLevelInfo {
		return
	}
	log.Output(2, fmt.Sprintf("I "+format, args...))
}

//WarnLog is the warn level log
func WarnLog(format string, args ...interface{}) {
	if logLevel < LogLevelWarn {
		return
	}
	log.Output(2, fmt.Sprintf("W "+format, args...))
}

//ErrorLog is the error level log
func ErrorLog(format string, args ...interface{}) {
	if logLevel < LogLevelError {
		return
	}
	log.Output(2, fmt.Sprintf("E "+format, args...))
}
