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
func Init() {
	var output io.Writer = os.Stderr

	// Configure output destination
	if logFile := os.Getenv("BACKLOG_LOG_FILE"); logFile != "" {
		// Ensure directory exists
		if err := os.MkdirAll(filepath.Dir(logFile), 0750); err != nil {
			// Fall back to stderr if we can't create the directory
			output = os.Stderr
		} else {
			file, err = os.OpenFile(logFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
			if err != nil {
				// Fall back to stderr if we can't open the log file
				output = os.Stderr
			} else {
				output = file
			}
		}
	}

	// Configure log level
	level := slog.LevelInfo
	if logLevel := os.Getenv("BACKLOG_LOG_LEVEL"); logLevel != "" {
		switch strings.ToLower(logLevel) {
		case "debug", "d":
			level = slog.LevelDebug
		case "info", "i":
			level = slog.LevelInfo
		case "warn", "warning", "w":
			level = slog.LevelWarn
		case "error", "err", "e":
			level = slog.LevelError
		}
	}

	// Configure format
	var handler slog.Handler
	if format := os.Getenv("BACKLOG_LOG_FORMAT"); strings.ToLower(format) == "json" || strings.ToLower(format) == "j" {
		handler = slog.NewJSONHandler(output, &slog.HandlerOptions{Level: level})
	} else {
		handler = slog.NewTextHandler(output, &slog.HandlerOptions{Level: level})
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
		Init()
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
