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

// Package cloudflare provides tests for the CloudflareProvider.
package cloudflare

import (
	"testing"
)

func TestCloudflareProvider_New_Unmarshal(t *testing.T) {
	cfgMap := map[string]any{
		"enabled":     true,
		"api_token":   "token",
		"zone_id":     "zone",
		"record_name": "name",
		"record_type": "A",
		"proxied":     true,
	}
	p, err := New(cfgMap)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Cfg.APIToken != "token" || p.Cfg.ZoneID != "zone" || !p.Cfg.Enabled {
		t.Errorf("config not unmarshaled correctly: %+v", p.Cfg)
	}
}

func TestCloudflareProvider_ProviderName(t *testing.T) {
	p := &CloudflareProvider{Cfg: &CloudflareConfig{}}
	if p.ProviderName() != "cloudflare" {
		t.Errorf("expected provider name 'cloudflare'")
	}
}
