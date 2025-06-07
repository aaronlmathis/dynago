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

// Package config provides configuration loading for dynago from YAML files.
//
// The configuration supports multiple DNS providers, logging options, and update intervals.
package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Config represents the root configuration for dynago loaded from YAML.
//
// Providers is a map of provider name to arbitrary config (for extensibility).
type Config struct {
	Interval  time.Duration  `yaml:"interval"`
	IPSource  string         `yaml:"ip_source"`
	LogLevel  string         `yaml:"log_level"`
	Providers map[string]any `yaml:"providers"`
}

// LoadConfig loads the configuration from the given YAML file path.
//
// It parses the YAML file, converts the interval string to time.Duration,
// and returns a Config struct or an error if parsing fails.
//
// Provider configs are left as generic maps for each provider.
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file %s: %w", path, err)
	}
	var raw struct {
		Interval  string         `yaml:"interval"`
		IPSource  string         `yaml:"ip_source"`
		LogLevel  string         `yaml:"log_level"`
		Providers map[string]any `yaml:"providers"`
	}
	if err := yaml.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("failed to parse config file %s: %w", path, err)
	}
	interval, err := time.ParseDuration(raw.Interval)
	if err != nil {
		return nil, fmt.Errorf("invalid interval %q in config file %s: %w", raw.Interval, path, err)
	}
	cfg := &Config{
		Interval:  interval,
		IPSource:  raw.IPSource,
		LogLevel:  raw.LogLevel,
		Providers: raw.Providers,
	}
	return cfg, nil
}

// ConfigFromMap parses a provider config from a generic map into a strongly-typed struct.
//
// This function is useful for converting provider-specific configuration
// from the generic map structure used in the YAML file to the
// strongly-typed struct expected by each provider's implementation.
//
// Usage:
//
//	var cfg CloudflareConfig
//	err := config.ConfigFromMap(rawMap, &cfg)
func ConfigFromMap(m any, out any) error {
	data, err := yaml.Marshal(m)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(data, out)
}
