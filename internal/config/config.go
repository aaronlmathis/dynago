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
// Fields:
//   - Interval: How often to check for IP changes (parsed as time.Duration).
//   - IPSource: URL to determine the public IP address.
//   - LogLevel: Logging level (debug, info, warn, error).
//   - Providers: DNS provider configurations.
type Config struct {
	Interval  time.Duration   `yaml:"interval"`
	IPSource  string          `yaml:"ip_source"`
	LogLevel  string          `yaml:"log_level"`
	Providers ProvidersConfig `yaml:"providers"`
}

// ProvidersConfig holds configuration for all supported DNS providers.
//
// Fields:
//   - Cloudflare: Cloudflare provider configuration.
//   - Route53: AWS Route53 provider configuration.
type ProvidersConfig struct {
	Cloudflare CloudflareConfig `yaml:"cloudflare"`
	Route53    Route53Config    `yaml:"route53"`
}

// CloudflareConfig holds Cloudflare-specific configuration.
//
// Fields:
//   - Enabled: Whether Cloudflare updates are enabled.
//   - APIToken: Cloudflare API token.
//   - ZoneID: Cloudflare zone ID.
//   - RecordName: DNS record name to update.
//   - RecordType: DNS record type (A or AAAA).
//   - Proxied: Whether the record should be proxied (orange cloud in Cloudflare UI).
type CloudflareConfig struct {
	Enabled    bool   `yaml:"enabled"`
	APIToken   string `yaml:"api_token"`
	ZoneID     string `yaml:"zone_id"`
	RecordName string `yaml:"record_name"`
	RecordType string `yaml:"record_type"`
	Proxied    bool   `yaml:"proxied"`
}

// Route53Config holds AWS Route53-specific configuration.
//
// Fields:
//   - Enabled: Whether Route53 updates are enabled.
//   - AccessKeyID: AWS access key ID.
//   - SecretAccessKey: AWS secret access key.
//   - HostedZoneID: AWS hosted zone ID.
//   - RecordName: DNS record name to update.
//   - RecordType: DNS record type (A or AAAA).
//   - Region: AWS region.
type Route53Config struct {
	Enabled         bool   `yaml:"enabled"`
	AccessKeyID     string `yaml:"access_key_id"`
	SecretAccessKey string `yaml:"secret_access_key"`
	HostedZoneID    string `yaml:"hosted_zone_id"`
	RecordName      string `yaml:"record_name"`
	RecordType      string `yaml:"record_type"`
	Region          string `yaml:"region"`
}

// LoadConfig loads the configuration from the given YAML file path.
//
// It parses the YAML file, converts the interval string to time.Duration,
// and returns a Config struct or an error if parsing fails.
//
// Example usage:
//
//	cfg, err := config.LoadConfig("configs/dynago.yml")
//	if err != nil {
//	    log.Fatal(err)
//	}
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file %s: %w", path, err)
	}
	var raw struct {
		Interval  string          `yaml:"interval"`
		IPSource  string          `yaml:"ip_source"`
		LogLevel  string          `yaml:"log_level"`
		Providers ProvidersConfig `yaml:"providers"`
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
