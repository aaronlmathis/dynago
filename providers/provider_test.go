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

// Package provider provides tests for the DNSProvider interface and registry.
package provider

import (
	"testing"
)

type mockProvider struct {
	name      string
	ip        string
	updateErr error
}

func (m *mockProvider) GetRecordIP() (string, error)   { return m.ip, nil }
func (m *mockProvider) UpdateRecordIP(ip string) error { m.ip = ip; return m.updateErr }
func (m *mockProvider) ProviderName() string           { return m.name }

func TestDNSProviderRegistry_AddsProviders(t *testing.T) {
	p1 := &mockProvider{name: "mock1", ip: "1.2.3.4"}
	p2 := &mockProvider{name: "mock2", ip: "2.3.4.5"}
	reg, err := NewDNSProviderRegistry(nil, p1, p2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(reg.Providers) != 2 {
		t.Errorf("expected 2 providers, got %d", len(reg.Providers))
	}
	if reg.Providers[0].ProviderName() != "mock1" || reg.Providers[1].ProviderName() != "mock2" {
		t.Errorf("provider names mismatch")
	}
}

func TestDNSProviderRegistry_Empty(t *testing.T) {
	_, err := NewDNSProviderRegistry(nil)
	if err == nil {
		t.Errorf("expected error when no providers enabled")
	}
}
