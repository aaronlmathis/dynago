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
// It supports different log levels, file output, and pretty console output in debug mode.

// Package logger provides tests for the logger package, ensuring correct logging behavior.
package logger

import (
	"bytes"
	"strings"
	"testing"

	"github.com/rs/zerolog"
)

// captureOutput is a helper function that captures log output written to appWriter.
// It sets the log level to DebugLevel, runs the provided function, and returns the captured output as a string.
func captureOutput(f func()) string {
	var buf bytes.Buffer
	appWriter = &buf
	logLevel = zerolog.DebugLevel
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	logger = zerolog.New(appWriter).With().Timestamp().Caller().Logger()
	f()
	return buf.String()
}

// TestInitLogger_Defaults verifies that InitLogger initializes without error when given an empty log file path and 'info' log level.
func TestInitLogger_Defaults(t *testing.T) {
	err := InitLogger("", "info")
	if err != nil {
		t.Fatalf("InitLogger failed: %v", err)
	}
}

// TestInfoLog checks that Info logs the expected formatted message to the application log.
func TestInfoLog(t *testing.T) {
	output := captureOutput(func() {
		Info("info message: %d", 42)
	})
	if !strings.Contains(output, "info message: 42") {
		t.Errorf("Info log not found in output: %s", output)
	}
}

// TestWarnLog checks that Warn logs the expected formatted warning message to the application log.
func TestWarnLog(t *testing.T) {
	output := captureOutput(func() {
		Warn("warn message: %s", "foo")
	})
	if !strings.Contains(output, "warn message: foo") {
		t.Errorf("Warn log not found in output: %s", output)
	}
}

// TestErrorLog checks that Error logs the expected formatted error message to the application log.
func TestErrorLog(t *testing.T) {
	output := captureOutput(func() {
		Error("error message: %v", "bar")
	})
	if !strings.Contains(output, "error message: bar") {
		t.Errorf("Error log not found in output: %s", output)
	}
}

// TestDebugLog_Enabled checks that Debug logs the expected message when logLevel is DebugLevel.
func TestDebugLog_Enabled(t *testing.T) {
	output := captureOutput(func() {
		Debug("debug message: %d", 99)
	})
	if !strings.Contains(output, "debug message: 99") {
		t.Errorf("Debug log not found in output: %s", output)
	}
}

// TestDebugLog_Disabled checks that Debug does not log anything when logLevel is higher than DebugLevel.
func TestDebugLog_Disabled(t *testing.T) {
	var buf bytes.Buffer
	appWriter = &buf
	logLevel = zerolog.InfoLevel
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	logger = zerolog.New(appWriter).With().Timestamp().Caller().Logger()
	Debug("should not appear")
	if buf.Len() != 0 {
		t.Errorf("Debug log written when disabled")
	}
}
