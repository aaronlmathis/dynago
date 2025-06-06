// dynago - Dynamic DNS updater for Cloudflare, Route 53, and more.
// Copyright (C) 2025  Aaron Mathis <aaron.mathis@gmail.com>
//
// This file is part of dynago.
//
// dynago is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// dynago is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with dynago.  If not, see <https://www.gnu.org/licenses/>.

// Package logger provides a simple logging system for the dynago application.
//
// It supports different log levels (debug, info, warn, error), writes logs to a file and pretty-prints to the console.
package logger

import (
	"io"
	"os"
	"strings"

	"github.com/rs/zerolog"
)

var (
	logger   zerolog.Logger // The global zerolog logger instance
	logLevel zerolog.Level  // Track the configured log level
)

// Writers for log output.
var (
	appWriter io.Writer // Writer for application logs (file or discard)
)

// InitLogger initializes the logging system for the application.
//
// All log messages are written to both the specified file and the console (with pretty formatting).
//
// Parameters:
//   - appLogFile:   Path to the application log file. If empty, logs are discarded.
//   - level:        Logging level ("debug", "info", "warn", "error").
//
// Returns an error if the log file cannot be opened.
//
// Example:
//
//	err := logger.InitLogger("app.log", "", "debug")
//	if err != nil {
//	    panic(err)
//	}
func InitLogger(appLogFile, level string) error {
	appWriter = io.Discard

	var appFile *os.File
	var err error

	if appLogFile != "" {
		appFile, err = os.OpenFile(appLogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		appWriter = appFile
	}

	switch strings.ToLower(level) {
	case "debug":
		logLevel = zerolog.DebugLevel
	case "info":
		logLevel = zerolog.InfoLevel
	case "warn":
		logLevel = zerolog.WarnLevel
	case "error":
		logLevel = zerolog.ErrorLevel
	default:
		logLevel = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(logLevel)

	// Set up pretty console output and multi-writer
	console := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: "2006-01-02 15:04:05"}
	multiWriter := io.MultiWriter(appWriter, console)
	logger = zerolog.New(multiWriter).With().Timestamp().Caller().Logger()
	return nil
}

// Info logs an informational message if the log level allows it.
//
//	format: Format string (like fmt.Printf).
//	args:   Arguments for the format string.
func Info(format string, args ...any) {
	if logLevel <= zerolog.InfoLevel {
		logger.Info().Msgf(format, args...)
	}
}

// Warn logs a warning message if the log level allows it.
//
//	format: Format string (like fmt.Printf).
//	args:   Arguments for the format string.
func Warn(format string, args ...any) {
	if logLevel <= zerolog.WarnLevel {
		logger.Warn().Msgf(format, args...)
	}
}

// Error logs an error message if the log level allows it.
//
//	format: Format string (like fmt.Printf).
//	args:   Arguments for the format string.
func Error(format string, args ...any) {
	if logLevel <= zerolog.ErrorLevel {
		logger.Error().Msgf(format, args...)
	}
}

// Debug logs a debug message if the log level allows it.
//
//	format: Format string (like fmt.Printf).
//	args:   Arguments for the format string.
func Debug(format string, args ...any) {
	if logLevel <= zerolog.DebugLevel {
		logger.Debug().Msgf(format, args...)
	}
}
