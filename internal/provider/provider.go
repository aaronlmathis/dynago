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

// Package provider implements DNS provider interfaces and concrete implementations for dynago.
//
// This package supports multiple DNS providers (Cloudflare, Route53) and provides a registry
// for enabled providers based on application configuration. Each provider implements the DNSProvider
// interface, which allows for querying and updating DNS records.
package provider

import (
	"context"
	"errors"
	"strings"

	"github.com/aaronlmathis/dynago/internal/config"
	"github.com/aaronlmathis/dynago/internal/logger"

	// AWS SDK imports
	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	route53types "github.com/aws/aws-sdk-go-v2/service/route53/types"

	// Cloudflare SDK import
	cloudflare "github.com/cloudflare/cloudflare-go"
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
//
// Returns a DNSProviderRegistry with all enabled providers, or an error if none are enabled.
func NewDNSProviderRegistry(cfg *config.Config) (*DNSProviderRegistry, error) {
	var providers []DNSProvider
	if cfg.Providers.Cloudflare.Enabled {
		providers = append(providers, &CloudflareProvider{cfg: &cfg.Providers.Cloudflare})
	}
	if cfg.Providers.Route53.Enabled {
		providers = append(providers, &Route53Provider{cfg: &cfg.Providers.Route53})
	}
	if len(providers) == 0 {
		return nil, errors.New("no DNS providers enabled in config")
	}
	return &DNSProviderRegistry{Providers: providers}, nil
}

// Route53Provider implements DNSProvider for AWS Route53.
//
// It uses the AWS SDK to query and update DNS records in a specified hosted zone.
type Route53Provider struct {
	cfg    *config.Route53Config // Provider-specific configuration
	client *route53.Client       // Cached AWS Route53 client
}

// ProviderName returns the string "route53" for AWS Route53 providers.
func (r *Route53Provider) ProviderName() string { return "route53" }

// getClient initializes and returns the AWS Route53 client, using static credentials from config.
func (r *Route53Provider) getClient(ctx context.Context) (*route53.Client, error) {
	if r.client != nil {
		return r.client, nil
	}
	awsCfg, err := awsconfig.LoadDefaultConfig(ctx,
		awsconfig.WithRegion(r.cfg.Region),
		awsconfig.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(r.cfg.AccessKeyID, r.cfg.SecretAccessKey, ""),
		),
	)
	if err != nil {
		return nil, err
	}
	r.client = route53.NewFromConfig(awsCfg)
	return r.client, nil
}

// GetRecordIP fetches the current IP address for the Route53 DNS record.
//
// Returns the IP address as a string, or an error if the record is not found or the API call fails.
func (r *Route53Provider) GetRecordIP() (string, error) {
	ctx := context.Background()
	client, err := r.getClient(ctx)
	if err != nil {
		return "", err
	}
	input := &route53.ListResourceRecordSetsInput{
		HostedZoneId:    aws.String(r.cfg.HostedZoneID),
		StartRecordName: aws.String(r.cfg.RecordName),
		StartRecordType: route53types.RRType(r.cfg.RecordType),
		MaxItems:        aws.Int32(1),
	}
	resp, err := client.ListResourceRecordSets(ctx, input)
	if err != nil {
		return "", err
	}
	for _, record := range resp.ResourceRecordSets {
		if strings.EqualFold(*record.Name, r.cfg.RecordName+".") && string(record.Type) == r.cfg.RecordType {
			if len(record.ResourceRecords) > 0 {
				return *record.ResourceRecords[0].Value, nil
			}
		}
	}
	return "", errors.New("record not found")
}

// UpdateRecordIP updates the Route53 DNS record to the given IP address.
//
// ip: The new IP address to set in the DNS record.
//
// Returns an error if the update fails.
func (r *Route53Provider) UpdateRecordIP(ip string) error {
	ctx := context.Background()
	client, err := r.getClient(ctx)
	if err != nil {
		return err
	}
	input := &route53.ChangeResourceRecordSetsInput{
		HostedZoneId: aws.String(r.cfg.HostedZoneID),
		ChangeBatch: &route53types.ChangeBatch{
			Changes: []route53types.Change{
				{
					Action: route53types.ChangeActionUpsert,
					ResourceRecordSet: &route53types.ResourceRecordSet{
						Name:            aws.String(r.cfg.RecordName),
						Type:            route53types.RRType(r.cfg.RecordType),
						TTL:             aws.Int64(300),
						ResourceRecords: []route53types.ResourceRecord{{Value: aws.String(ip)}},
					},
				},
			},
		},
	}
	_, err = client.ChangeResourceRecordSets(ctx, input)
	logger.Info("Updating Route53 record %s to IP %s", r.cfg.RecordName, ip)

	return err
}

// CloudflareProvider implements DNSProvider for Cloudflare.
//
// It uses the Cloudflare Go SDK to query and update DNS records in a specified zone.
type CloudflareProvider struct {
	cfg    *config.CloudflareConfig // Provider-specific configuration
	client *cloudflare.API          // Cached Cloudflare API client
}

// ProviderName returns the string "cloudflare" for Cloudflare providers.
func (c *CloudflareProvider) ProviderName() string { return "cloudflare" }

// getClient initializes and returns the Cloudflare API client using the API token from config.
func (c *CloudflareProvider) getClient() (*cloudflare.API, error) {
	if c.client != nil {
		return c.client, nil
	}
	api, err := cloudflare.NewWithAPIToken(c.cfg.APIToken)
	if err != nil {
		return nil, err
	}
	c.client = api
	return c.client, nil
}

// GetRecordIP fetches the current IP address for the Cloudflare DNS record.
//
// Returns the IP address as a string, or an error if the record is not found or the API call fails.
func (c *CloudflareProvider) GetRecordIP() (string, error) {
	client, err := c.getClient()
	if err != nil {
		return "", err
	}
	zone := cloudflare.ZoneIdentifier(c.cfg.ZoneID)
	records, _, err := client.ListDNSRecords(context.Background(), zone, cloudflare.ListDNSRecordsParams{
		Name: c.cfg.RecordName,
		Type: c.cfg.RecordType,
	})
	if err != nil {
		return "", err
	}
	for _, record := range records {
		if record.Name == c.cfg.RecordName && record.Type == c.cfg.RecordType {
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
	zone := cloudflare.ZoneIdentifier(c.cfg.ZoneID)
	records, _, err := client.ListDNSRecords(context.Background(), zone, cloudflare.ListDNSRecordsParams{
		Name: c.cfg.RecordName,
		Type: c.cfg.RecordType,
	})
	if err != nil {
		return err
	}
	for _, record := range records {
		if record.Name == c.cfg.RecordName && record.Type == c.cfg.RecordType {
			edit := cloudflare.UpdateDNSRecordParams{
				ID:      record.ID,
				Type:    c.cfg.RecordType,
				Name:    c.cfg.RecordName,
				Content: ip,
				Proxied: &c.cfg.Proxied,
			}
			logger.Info("Updating Cloudflare record %s (ID: %s, zone: %s) to IP %s", c.cfg.RecordName, record.ID, c.cfg.ZoneID, ip)
			_, err := client.UpdateDNSRecord(context.Background(), zone, edit)
			return err
		}
	}
	logger.Error("Cloudflare: record not found for update. Name: %s, Type: %s, Zone: %s. Records returned: %d", c.cfg.RecordName, c.cfg.RecordType, c.cfg.ZoneID, len(records))
	return errors.New("record not found for update")
}
