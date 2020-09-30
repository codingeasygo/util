package socks

import "testing"

func TestLog(t *testing.T) {
	//
	SetLogLevel(LogLevelDebug)
	DebugLog("debug")
	InfoLog("info")
	WarnLog("warn")
	ErrorLog("error")
	//
	SetLogLevel(LogLevelInfo)
	DebugLog("debug")
	InfoLog("info")
	WarnLog("warn")
	ErrorLog("error")
	//
	SetLogLevel(LogLevelWarn)
	DebugLog("debug")
	InfoLog("info")
	WarnLog("warn")
	ErrorLog("error")
	//
	SetLogLevel(LogLevelError)
	DebugLog("debug")
	InfoLog("info")
	WarnLog("warn")
	ErrorLog("error")
	//
	SetLogLevel(1)
	ErrorLog("error")
}
