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

// Package provider defines the DNSProvider interface and registry for dynago.
//
// To add a new provider, implement the DNSProvider interface in your own package.
// Your provider should define its own config struct and unmarshal from a generic map
// using config.ConfigFromMap. See the cloudflare and route53 packages for examples.
//
// At runtime, the main config file's `providers:` section is a map, and each provider
// receives its own sub-map. The provider is responsible for parsing and validating its config.
package provider

import (
	"errors"

	"github.com/aaronlmathis/dynago/internal/config"
)

// DNSProvider defines the interface for DNS providers.
//
// Implementations must provide methods to get and update the DNS record IP.
type DNSProvider interface {
	// GetRecordIP returns the current IP address configured in the DNS record.
	GetRecordIP() (string, error)
	// UpdateRecordIP updates the DNS record to the given IP address.
	UpdateRecordIP(ip string) error
	// ProviderName returns the name of the provider (e.g., "cloudflare", "route53").
	ProviderName() string
}

// DNSProviderRegistry holds all enabled DNS providers.
//
// Providers is a slice of DNSProvider implementations that are enabled in the config.
type DNSProviderRegistry struct {
	Providers []DNSProvider
}

// NewDNSProviderRegistry creates a registry of enabled DNS providers based on config.
//
// cfg: The loaded application configuration.
// providers: One or more DNSProvider implementations to register.
//
// Returns a DNSProviderRegistry with all enabled providers, or an error if none are enabled.
func NewDNSProviderRegistry(cfg *config.Config, providers ...DNSProvider) (*DNSProviderRegistry, error) {
	if len(providers) == 0 {
		return nil, errors.New("no DNS providers enabled in config")
	}
	return &DNSProviderRegistry{Providers: providers}, nil
}
