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

// Package route53 provides tests for the Route53Provider.
package route53

import (
	"testing"
)

func TestRoute53Provider_New_Unmarshal(t *testing.T) {
	cfgMap := map[string]any{
		"enabled":           true,
		"access_key_id":     "id",
		"secret_access_key": "secret",
		"hosted_zone_id":    "zone",
		"record_name":       "name",
		"record_type":       "A",
		"region":            "us-east-1",
	}
	p, err := New(cfgMap)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Cfg.AccessKeyID != "id" || p.Cfg.HostedZoneID != "zone" || !p.Cfg.Enabled {
		t.Errorf("config not unmarshaled correctly: %+v", p.Cfg)
	}
}

func TestRoute53Provider_ProviderName(t *testing.T) {
	p := &Route53Provider{Cfg: &Route53Config{}}
	if p.ProviderName() != "route53" {
		t.Errorf("expected provider name 'route53'")
	}
}
