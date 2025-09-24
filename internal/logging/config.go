// Package logging provides global logging configuration and utilities
// for the backlog application.
package logging

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

var (
	logger *slog.Logger
	file   *os.File
)

// Init initializes the global logger based on environment variables
// BACKLOG_LOG_FILE: path to log file (optional, defaults to stderr)
// BACKLOG_LOG_LEVEL: log level (debug, info, warn, error, defaults to info)
// BACKLOG_LOG_FORMAT: format (json, text, defaults to text)
func Init(level, format, logFile string) {
	var output io.Writer = os.Stderr

	// Configure output destination
	if logFile != "" {
		// Ensure directory exists
		if err := os.MkdirAll(filepath.Dir(logFile), 0o750); err != nil {
			// Fall back to stderr if we can't create the directory
			output = os.Stderr
		} else {
			file, err = os.OpenFile(logFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o600)
			if err != nil {
				// Fall back to stderr if we can't open the log file
				output = os.Stderr
			} else {
				output = file
			}
		}
	}

	// Configure log level
	lvl := slog.LevelInfo
	switch strings.ToLower(level) {
	case "debug", "d":
		lvl = slog.LevelDebug
	case "info", "i":
		lvl = slog.LevelInfo
	case "warn", "warning", "w":
		lvl = slog.LevelWarn
	case "error", "err", "e":
		lvl = slog.LevelError
	}

	opts := &slog.HandlerOptions{
		Level:     lvl,
		AddSource: false,
		ReplaceAttr: func(_ []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				t := a.Value.Time()
				a = slog.Attr{Key: "time", Value: slog.StringValue(t.Format("15:04:05"))}
			}
			return a
		},
	}
	// Configure format
	var handler slog.Handler
	if strings.ToLower(format) == "json" || strings.ToLower(format) == "j" {
		handler = slog.NewJSONHandler(output, opts)
	} else {
		handler = slog.NewTextHandler(output, opts)
	}

	logger = slog.New(handler)
	slog.SetDefault(logger)
}

func Close() {
	if file != nil {
		if err := file.Close(); err != nil {
			fmt.Printf("closing file: %v", err)
		}
	}
}

// GetLogger returns the configured logger instance
func GetLogger() *slog.Logger {
	if logger == nil {
		Init("debug", "text", "")
	}
	return logger
}

// Debug logs a debug message
func Debug(msg string, args ...any) {
	GetLogger().Debug(msg, args...)
}

// Info logs an info message
func Info(msg string, args ...any) {
	GetLogger().Info(msg, args...)
}

// Warn logs a warning message
func Warn(msg string, args ...any) {
	GetLogger().Warn(msg, args...)
}

// Error logs an error message
func Error(msg string, args ...any) {
	GetLogger().Error(msg, args...)
}
