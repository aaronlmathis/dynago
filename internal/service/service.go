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

// Package service provides the DNS update service for dynago.
//
// The DNSUpdateService periodically checks the current public IP address and updates DNS records
// via enabled providers (Cloudflare, Route53, etc.) if the IP has changed.
package service

import (
	"context"
	"fmt"
	"time"

	"github.com/aaronlmathis/dynago/internal/config"
	"github.com/aaronlmathis/dynago/internal/logger"
	"github.com/aaronlmathis/dynago/internal/utils"
	providers "github.com/aaronlmathis/dynago/providers"
	cfprovider "github.com/aaronlmathis/dynago/providers/cloudflare"
	r53provider "github.com/aaronlmathis/dynago/providers/route53"
)

// DNSUpdateService manages the periodic update of DNS records for the host's current public IP.
//
// It loads configuration, initializes providers, and runs a loop to check and update DNS records as needed.
type DNSUpdateService struct {
	cfg      *config.Config        // Application configuration
	ctx      context.Context       // Service context for cancellation
	provider providers.DNSProvider // (Unused, reserved for future single-provider mode)
}

// NewDNSUpdateService creates a new DNSUpdateService with the given context and configuration.
//
// ctx: Context for cancellation and shutdown.
// cfg: Loaded application configuration.
func NewDNSUpdateService(ctx context.Context, cfg *config.Config) *DNSUpdateService {
	return &DNSUpdateService{
		cfg: cfg,
		ctx: ctx,
	}
}

// Start begins the DNS update loop.
//
// It periodically fetches the current public IP address using the configured source,
// compares it to the DNS records for each enabled provider, and updates the records if the IP has changed.
//
// Returns an error if the service cannot start or if no providers are enabled.
func (s *DNSUpdateService) Start() error {
	if s.cfg == nil {
		return nil // No configuration provided, nothing to do.
	}
	logger.Info("DNSUpdateService starting... ")

	var providersList []providers.DNSProvider
	if raw, ok := s.cfg.Providers["cloudflare"]; ok {
		cf, err := cfprovider.New(raw)
		if err == nil && cf.Cfg.Enabled {
			providersList = append(providersList, cf)
		}
	}
	if raw, ok := s.cfg.Providers["route53"]; ok {
		r53, err := r53provider.New(raw)
		if err == nil && r53.Cfg.Enabled {
			providersList = append(providersList, r53)
		}
	}
	reg, err := providers.NewDNSProviderRegistry(s.cfg, providersList...)
	if err != nil {
		logger.Error("No DNS providers enabled: %v", err)
		return fmt.Errorf("failed to create DNS provider registry: %w", err)
	}

	interval := s.cfg.Interval
	ipSource := s.cfg.IPSource
	lastKnown := make(map[string]string) // providerName -> last IP

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-s.ctx.Done():
			logger.Info("DNSUpdateService stopped")
			return nil
		case <-ticker.C:
			currentIP, err := utils.GetCurrentIP(ipSource)
			if err != nil {
				logger.Error("Failed to get current IP: %v", err)
				continue
			}
			for _, p := range reg.Providers {
				providerName := p.ProviderName()
				dnsIP, err := p.GetRecordIP()
				if err != nil {
					logger.Error("%s: failed to get DNS record IP: %v", providerName, err)
					continue
				}
				if dnsIP != currentIP {
					logger.Info("%s: IP mismatch (current: %s, DNS: %s), updating...", providerName, currentIP, dnsIP)
					if err := p.UpdateRecordIP(currentIP); err != nil {
						logger.Error("%s: failed to update DNS record: %v", providerName, err)
						continue
					}
					lastKnown[providerName] = currentIP
					logger.Info("%s: DNS record updated to %s", providerName, currentIP)
				} else {
					logger.Debug("%s: IP unchanged (%s)", providerName, currentIP)
				}
			}
		}
	}
}

// Stop stops the DNS update service and performs any necessary cleanup.
//
// Returns an error if cleanup fails (currently always returns nil).
func (s *DNSUpdateService) Stop() error {
	// Clean up resources, close connections, etc.
	// This is a placeholder for actual cleanup logic.
	return nil
}
