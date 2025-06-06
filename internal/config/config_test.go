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

// Package config provides tests for the configuration loader.
package config

import (
	"os"
	"testing"
)

const sampleYAML = `
interval: 5m
ip_source: "https://api.ipify.org"
log_level: "info"
providers:
  cloudflare:
    enabled: true
    api_token: "cf-token"
    zone_id: "cf-zone"
    record_name: "home.example.com"
    record_type: "A"
    proxied: true
  route53:
    enabled: false
    access_key_id: "aws-key"
    secret_access_key: "aws-secret"
    hosted_zone_id: "aws-zone"
    record_name: "home.example.com"
    record_type: "A"
    region: "us-east-1"
`

// TestLoadConfig checks that LoadConfig correctly parses a sample YAML config file.
func TestLoadConfig(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "dynago-config-*.yml")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpfile.Name())
	if _, err := tmpfile.Write([]byte(sampleYAML)); err != nil {
		t.Fatalf("failed to write sample YAML: %v", err)
	}
	tmpfile.Close()

	cfg, err := LoadConfig(tmpfile.Name())
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}
	if cfg.Interval.Minutes() != 5 {
		t.Errorf("expected interval 5m, got %v", cfg.Interval)
	}
	if cfg.IPSource != "https://api.ipify.org" {
		t.Errorf("unexpected ip_source: %s", cfg.IPSource)
	}
	if cfg.LogLevel != "info" {
		t.Errorf("unexpected log_level: %s", cfg.LogLevel)
	}
	if !cfg.Providers.Cloudflare.Enabled {
		t.Errorf("cloudflare.enabled should be true")
	}
	if cfg.Providers.Cloudflare.APIToken != "cf-token" {
		t.Errorf("unexpected cloudflare.api_token: %s", cfg.Providers.Cloudflare.APIToken)
	}
	if cfg.Providers.Cloudflare.Proxied != true {
		t.Errorf("cloudflare.proxied should be true")
	}
	if cfg.Providers.Route53.Enabled {
		t.Errorf("route53.enabled should be false")
	}
	if cfg.Providers.Route53.Region != "us-east-1" {
		t.Errorf("unexpected route53.region: %s", cfg.Providers.Route53.Region)
	}
}
