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

package utils

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestGetCurrentIP checks that GetCurrentIP fetches the IP from a mock HTTP server.
func TestGetCurrentIP(t *testing.T) {
	mockIP := "1.2.3.4"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(mockIP))
	}))
	defer ts.Close()

	ip, err := GetCurrentIP(ts.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ip != mockIP {
		t.Errorf("expected %s, got %s", mockIP, ip)
	}
}
