package util

import (
	"fmt"
	"go.uber.org/zap/zapcore"
	"log"
	"os"
)

var _Log *log.Logger

func init() {
	_Log = log.New(os.Stdout, "[Console] ", log.LstdFlags)
}

var logLevel = zapcore.InfoLevel

func consoleTag(lv zapcore.Level, tag, format string, args ...interface{}) {
	if lv >= logLevel {
		_Log.Printf("[%s] %s\n", tag, fmt.Sprintf(format, args...))
	}
}

func SetLevel(level zapcore.Level) {
	logLevel = level
}

func Debug(format string, args ...interface{}) {
	consoleTag(zapcore.DebugLevel, "DEBUG", format, args...)
}

func Info(format string, args ...interface{}) {
	consoleTag(zapcore.InfoLevel, "INFO", format, args...)
}

func Warn(format string, args ...interface{}) {
	consoleTag(zapcore.WarnLevel, "WARN", format, args...)
}

func Error(format string, args ...interface{}) {
	consoleTag(zapcore.ErrorLevel, "ERROR", format, args...)
}

func Panic(format string, args ...interface{}) {
	_Log.Panicf("[%s] %s\n", "FALTA", fmt.Sprintf(format, args...))
}
