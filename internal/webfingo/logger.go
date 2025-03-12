package webfingo

import (
	"io"
	"log"
	"os"
)

// Logger provides a simple logging interface that can be silenced for tests
type Logger struct {
	*log.Logger
}

// DefaultLogger is the default logger instance used by the application
var DefaultLogger = NewLogger(os.Stderr)

// NewLogger creates a new Logger that writes to the specified writer
func NewLogger(w io.Writer) *Logger {
	return &Logger{
		Logger: log.New(w, "", log.LstdFlags),
	}
}

// SilentLogger returns a logger that doesn't output anything
func SilentLogger() *Logger {
	return NewLogger(io.Discard)
}

// Fatalf is a wrapper for Logger.Fatalf
func (l *Logger) Fatalf(format string, v ...interface{}) {
	l.Logger.Fatalf(format, v...)
}

// Printf is a wrapper for Logger.Printf
func (l *Logger) Printf(format string, v ...interface{}) {
	l.Logger.Printf(format, v...)
}

// Println is a wrapper for Logger.Println
func (l *Logger) Println(v ...interface{}) {
	l.Logger.Println(v...)
}
