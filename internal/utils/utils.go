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

// Package utils provides utility functions for dynago, such as public IP detection.
//
// These helpers are used by the service layer to determine the current public IP address of the host.
package utils

import (
	"errors"
	"io"
	"net/http"
)

// GetCurrentIP fetches the current public IP address from the specified source URL.
//
// ipSource: The URL of an external service that returns the public IP as plain text (e.g., https://api.ipify.org).
//
// Returns the IP address as a string, or an error if the request fails or the response is not HTTP 200.
//
// Example:
//
//	ip, err := utils.GetCurrentIP("https://api.ipify.org")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println("Current IP:", ip)
func GetCurrentIP(ipSource string) (string, error) {
	resp, err := http.Get(ipSource)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", errors.New("failed to fetch IP: non-200 response")
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}
