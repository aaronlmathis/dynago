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

// Package cloudflare implements the DNSProvider interface for Cloudflare.
//
// This package provides a CloudflareProvider type that can be registered with the dynago DNS update service.
// It uses the Cloudflare Go SDK to query and update DNS records in a specified zone.
package cloudflare

import (
	"context"
	"errors"

	"github.com/aaronlmathis/dynago/internal/config"
	cf "github.com/cloudflare/cloudflare-go"
)

// CloudflareConfig holds Cloudflare-specific configuration.
type CloudflareConfig struct {
	Enabled    bool   `yaml:"enabled"`
	APIToken   string `yaml:"api_token"`
	ZoneID     string `yaml:"zone_id"`
	RecordName string `yaml:"record_name"`
	RecordType string `yaml:"record_type"`
	Proxied    bool   `yaml:"proxied"`
}

// CloudflareProvider implements the DNSProvider interface for Cloudflare.
//
// It uses the Cloudflare Go SDK to query and update DNS records in a specified zone.
type CloudflareProvider struct {
	Cfg    *CloudflareConfig // Provider-specific configuration
	Client *cf.API           // Cached Cloudflare API client
}

// New creates a new CloudflareProvider from a generic config map.
//
// Usage: cfprovider.New(configMap)
func New(raw any) (*CloudflareProvider, error) {
	var cfg CloudflareConfig
	err := config.ConfigFromMap(raw, &cfg)
	if err != nil {
		return nil, err
	}
	return &CloudflareProvider{Cfg: &cfg}, nil
}

// getClient initializes and returns the Cloudflare API client using the API token from config.
func (c *CloudflareProvider) getClient() (*cf.API, error) {
	if c.Client != nil {
		return c.Client, nil
	}
	api, err := cf.NewWithAPIToken(c.Cfg.APIToken)
	if err != nil {
		return nil, err
	}
	c.Client = api
	return c.Client, nil
}

// ProviderName returns the string "cloudflare" for Cloudflare providers.
func (c *CloudflareProvider) ProviderName() string { return "cloudflare" }

// GetRecordIP fetches the current IP address for the Cloudflare DNS record.
//
// Returns the IP address as a string, or an error if the record is not found or the API call fails.
func (c *CloudflareProvider) GetRecordIP() (string, error) {
	client, err := c.getClient()
	if err != nil {
		return "", err
	}
	zone := cf.ZoneIdentifier(c.Cfg.ZoneID)
	records, _, err := client.ListDNSRecords(context.Background(), zone, cf.ListDNSRecordsParams{
		Name: c.Cfg.RecordName,
		Type: c.Cfg.RecordType,
	})
	if err != nil {
		return "", err
	}
	for _, record := range records {
		if record.Name == c.Cfg.RecordName && record.Type == c.Cfg.RecordType {
			return record.Content, nil
		}
	}
	return "", errors.New("record not found")
}

// UpdateRecordIP updates the Cloudflare DNS record to the given IP address.
//
// ip: The new IP address to set in the DNS record. The proxied status is set according to config.
//
// Returns an error if the update fails or the record is not found.
func (c *CloudflareProvider) UpdateRecordIP(ip string) error {
	client, err := c.getClient()
	if err != nil {
		return err
	}
	zone := cf.ZoneIdentifier(c.Cfg.ZoneID)
	records, _, err := client.ListDNSRecords(context.Background(), zone, cf.ListDNSRecordsParams{
		Name: c.Cfg.RecordName,
		Type: c.Cfg.RecordType,
	})
	if err != nil {
		return err
	}
	for _, record := range records {
		if record.Name == c.Cfg.RecordName && record.Type == c.Cfg.RecordType {
			edit := cf.UpdateDNSRecordParams{
				ID:      record.ID,
				Type:    c.Cfg.RecordType,
				Name:    c.Cfg.RecordName,
				Content: ip,
				Proxied: &c.Cfg.Proxied,
			}
			_, err := client.UpdateDNSRecord(context.Background(), zone, edit)
			return err
		}
	}
	return errors.New("record not found for update")
}
