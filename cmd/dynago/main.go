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

// Command dynago is the entry point for the dynago dynamic DNS updater application.
//
// This command-line tool loads configuration, initializes logging, and starts the DNS update service.
// It supports updating DNS records for multiple providers (Cloudflare, Route53, etc.) based on the
// current public IP address of the host. The service runs until interrupted (SIGINT/SIGTERM).
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/aaronlmathis/dynago/internal/config"
	"github.com/aaronlmathis/dynago/internal/logger"
	"github.com/aaronlmathis/dynago/internal/service"
)

var (
	Version    = "dev"     // Application version (set at build time)
	BuildTime  = "unknown" // Build timestamp (set at build time)
	GitCommit  = "none"    // Git commit hash (set at build time)
	ConfigPath string      // Path to the configuration file
	LogFile    string      // Path to the log file (optional)
)

// run loads configuration, initializes logging, and starts the DNS update service.
//
// Returns an error if configuration or logger initialization fails, or if the service fails to start.
func run() error {
	// Load the configuration from the config file.
	cfg, err := config.LoadConfig(ConfigPath)
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Initialize the logger with the configured log level
	// and log file path from the configuration.
	if err := logger.InitLogger(LogFile, cfg.LogLevel); err != nil {
		return fmt.Errorf("failed to initialize logger: %w", err)

	}

	logger.Debug("Starting dynago version %s (built at %s, commit %s)", Version, BuildTime, GitCommit)
	logger.Debug("Configuration loaded from %s", ConfigPath)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
		<-sig
		cancel()
	}()

	// Create the DNS update service with the loaded configuration.
	dnsService := service.NewDNSUpdateService(ctx, cfg)
	if err := dnsService.Start(); err != nil {
		return fmt.Errorf("failed to start DNS update service: %w", err)
	}

	return nil

}

// main is the entry point for the dynago application.
//
// It parses command-line flags, then calls run(). If an error occurs, it prints the error and exits with status 1.
func main() {
	versionFlag := flag.Bool("version", false, "print version information and exit")
	if *versionFlag {
		fmt.Printf("dynago %s (commit %s, built %s)\n", Version, GitCommit, BuildTime)
		os.Exit(0)
	}

	flag.StringVar(&ConfigPath, "config", "configs/dynago.yml", "Path to the configuration file")
	flag.StringVar(&LogFile, "log", "", "Path to the log file (optional, defaults to stdout)")
	flag.Parse()
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
