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

package provider

import (
	"testing"

	"github.com/aaronlmathis/dynago/internal/config"
)

// TestDNSProviderRegistry ensures the registry enables only configured providers.
func TestDNSProviderRegistry(t *testing.T) {
	cfg := &config.Config{
		Providers: config.ProvidersConfig{
			Cloudflare: config.CloudflareConfig{Enabled: true},
			Route53:    config.Route53Config{Enabled: false},
		},
	}
	reg, err := NewDNSProviderRegistry(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(reg.Providers) != 1 || reg.Providers[0].ProviderName() != "cloudflare" {
		t.Errorf("expected only cloudflare provider enabled")
	}

	cfg.Providers.Cloudflare.Enabled = false
	cfg.Providers.Route53.Enabled = true
	reg, err = NewDNSProviderRegistry(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(reg.Providers) != 1 || reg.Providers[0].ProviderName() != "route53" {
		t.Errorf("expected only route53 provider enabled")
	}

	cfg.Providers.Cloudflare.Enabled = false
	cfg.Providers.Route53.Enabled = false
	_, err = NewDNSProviderRegistry(cfg)
	if err == nil {
		t.Errorf("expected error when no providers enabled")
	}
}

// TestProviderStubs ensures stub methods return expected values.
func TestProviderStubs(t *testing.T) {
	cf := &CloudflareProvider{cfg: &config.CloudflareConfig{APIToken: "dummy", ZoneID: "dummy", RecordName: "dummy", RecordType: "A"}}
	if _, err := cf.GetRecordIP(); err == nil {
		t.Errorf("expected error from GetRecordIP with dummy config")
	}
	if err := cf.UpdateRecordIP("1.2.3.4"); err == nil {
		t.Errorf("expected error from UpdateRecordIP with dummy config")
	}

	r53 := &Route53Provider{cfg: &config.Route53Config{AccessKeyID: "dummy", SecretAccessKey: "dummy", HostedZoneID: "dummy", RecordName: "dummy", RecordType: "A", Region: "us-east-1"}}
	if _, err := r53.GetRecordIP(); err == nil {
		t.Errorf("expected error from GetRecordIP with dummy config")
	}
	if err := r53.UpdateRecordIP("1.2.3.4"); err == nil {
		t.Errorf("expected error from UpdateRecordIP with dummy config")
	}
}
