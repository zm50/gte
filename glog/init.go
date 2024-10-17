package glog

import (
	"github.com/go75/gte/gconf"
	"github.com/go75/gte/trait"

	"github.com/natefinch/lumberjack"
)

var log trait.Log

// Init initializes the log client.
func Init() {
	out := &lumberjack.Logger{
		Filename:   gconf.Config.LogFilename(),      // 日志文件存放目录
		MaxSize:    gconf.Config.LogMaxSize(),       // 文件大小限制,单位MB
		MaxBackups: gconf.Config.LogMaxBackups(),    // 最大保留日志文件数量
		MaxAge:     gconf.Config.LogMaxAge(),        // 日志文件保留天数
		Compress:   gconf.Config.LogCompress(),      // 是否压缩处理
	}

	log = NewLog(out)
}

// Info logs a message with level Info.
func Info(args ...any) {
	log.Info(args...)
}

// Error logs a message with level Error.
func Error(args ...any) {
	log.Error(args...)
}

// Warn logs a message with level Warn.
func Warn(args ...any) {
	log.Warn(args...)
}

// Fatal logs a message with level Fatal.
func Fatal(args ...any) {
	log.Fatal(args...)
}

// Infof logs a message with level Info.
func Infof(format string, args ...any) {
	log.Infof(format, args...)
}

// Errorf logs a message with level Error.
func Errorf(format string, args ...any) {
	log.Errorf(format, args...)
}

// Warnf logs a message with level Warn.
func Warnf(format string, args ...any) {
	log.Warnf(format, args...)
}

// Fatalf logs a message with level Fatal.
func Fatalf(format string, args ...any) {
	log.Fatalf(format, args...)
}