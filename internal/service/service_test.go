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

package service

import (
	"context"
	"testing"
	"time"

	"github.com/aaronlmathis/dynago/internal/config"
)

type mockProvider struct {
	name      string
	getIP     string
	updatedIP string
	getErr    error
	updateErr error
}

func (m *mockProvider) GetRecordIP() (string, error)   { return m.getIP, m.getErr }
func (m *mockProvider) UpdateRecordIP(ip string) error { m.updatedIP = ip; return m.updateErr }
func (m *mockProvider) ProviderName() string           { return m.name }

func TestDNSUpdateService_Start(t *testing.T) {
	cfg := &config.Config{Interval: 10 * time.Millisecond, IPSource: "mock", LogLevel: "debug"}
	ctx, cancel := context.WithCancel(context.Background())
	service := NewDNSUpdateService(ctx, cfg)

	// Patch provider registry and utils for test
	mockProv := &mockProvider{name: "mock", getIP: "4.3.2.1"}
	_ = mockProv
	_ = service
	// TODO: Refactor service for better testability (dependency injection)

	// Simulate one tick then cancel
	go func() {
		time.Sleep(20 * time.Millisecond)
		cancel()
	}()
}
