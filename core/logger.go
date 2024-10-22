package core

import (
	"github.com/natefinch/lumberjack"
	"github.com/zm50/gte/trait"
	"k8s.io/klog/v2"
)

// Logger is a struct for logging
type Logger struct {}

var _ trait.Logger = (*Logger)(nil)

// Init creates a new Log object.
func NewLogger(filename string, maxSize int, maxBackups int, maxAge int, compress bool) trait.Logger {
	out := &lumberjack.Logger{
		Filename:   filename,   // 日志文件存放目录
		MaxSize:    maxSize,    // 文件大小限制,单位MB
		MaxBackups: maxBackups, // 最大保留日志文件数量
		MaxAge:     maxAge,     // 日志文件保留天数
		Compress:   compress,   // 是否压缩处理
	}

	klog.SetOutput(out)

	return &Logger{}
}

// Info logs a message at info level
func (l *Logger) Info(args ...any) {
	klog.InfolnDepth(2, args...)
}

// Error logs a message at error level
func (l *Logger) Error(args ...any) {
	klog.ErrorlnDepth(2, args...)
}

// Warn logs a message at warning level
func (l *Logger) Warn(args ...any) {
	klog.WarninglnDepth(2, args...)
}

// Fatal logs a message at fatal level
func (l *Logger) Fatal(args ...any) {
	klog.FatallnDepth(2, args...)
}

// Infof logs a message at info level with format
func (l *Logger) Infof(format string, args ...any) {
	klog.InfofDepth(2, format, args...)
}

// Errorf logs a message at error level with format
func (l *Logger) Errorf(format string, args ...any) {
	klog.ErrorfDepth(2, format, args...)
}

// Warnf logs a message at warning level with format
func (l *Logger) Warnf(format string, args ...any) {
	klog.WarningfDepth(2, format, args...)
}

// Fatal logs a message at fatal level with format
func (l *Logger) Fatalf(format string, args ...any) {
	klog.FatalfDepth(2, format, args...)
}
