package glog

import (
	"io"

	"k8s.io/klog/v2"
)

// Log is a struct for logging
type Log struct {}

// NewLog creates a new Log object
func NewLog(output io.Writer) *Log {
	klog.SetOutput(output)
	return &Log{}
}

// Info logs a message at info level
func (l *Log) Info(args ...any) {
	klog.InfolnDepth(2, args...)
}

// Error logs a message at error level
func (l *Log) Error(args ...any) {
	klog.ErrorlnDepth(2, args...)
}

// Warn logs a message at warning level
func (l *Log) Warn(args ...any) {
	klog.WarninglnDepth(2, args...)
}

// Fatal logs a message at fatal level
func (l *Log) Fatal(args ...any) {
	klog.FatallnDepth(2, args...)
}

// Infof logs a message at info level with format
func (l *Log) Infof(format string, args ...any) {
	klog.InfofDepth(2, format, args...)
}

// Errorf logs a message at error level with format
func (l *Log) Errorf(format string, args ...any) {
	klog.ErrorfDepth(2, format, args...)
}

// Warnf logs a message at warning level with format
func (l *Log) Warnf(format string, args ...any) {
	klog.WarningfDepth(2, format, args...)
}

// Fatal logs a message at fatal level with format
func (l *Log) Fatalf(format string, args ...any) {
	klog.FatalfDepth(2, format, args...)
}