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

// Package route53 implements the DNSProvider interface for AWS Route53.
//
// This package provides a Route53Provider type that can be registered with the dynago DNS update service.
// It uses the AWS SDK to query and update DNS records in a specified hosted zone.
package route53

import (
	"context"
	"errors"
	"strings"

	"github.com/aaronlmathis/dynago/internal/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	r53types "github.com/aws/aws-sdk-go-v2/service/route53/types"
)

// Route53Config holds AWS Route53-specific configuration.
type Route53Config struct {
	Enabled         bool   `yaml:"enabled"`
	AccessKeyID     string `yaml:"access_key_id"`
	SecretAccessKey string `yaml:"secret_access_key"`
	HostedZoneID    string `yaml:"hosted_zone_id"`
	RecordName      string `yaml:"record_name"`
	RecordType      string `yaml:"record_type"`
	Region          string `yaml:"region"`
}

// Route53Provider implements the DNSProvider interface for AWS Route53.
//
// It uses the AWS SDK to query and update DNS records in a specified hosted zone.
type Route53Provider struct {
	Cfg    *Route53Config  // Provider-specific configuration
	Client *route53.Client // Cached AWS Route53 client
}

// New creates a new Route53Provider from a generic config map.
//
// Usage: route53provider.New(configMap)
func New(raw any) (*Route53Provider, error) {
	var cfg Route53Config
	err := config.ConfigFromMap(raw, &cfg)
	if err != nil {
		return nil, err
	}
	return &Route53Provider{Cfg: &cfg}, nil
}

// getClient initializes and returns the AWS Route53 client, using static credentials from config.
func (r *Route53Provider) getClient(ctx context.Context) (*route53.Client, error) {
	if r.Client != nil {
		return r.Client, nil
	}
	awsCfg, err := awsconfig.LoadDefaultConfig(ctx,
		awsconfig.WithRegion(r.Cfg.Region),
		awsconfig.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(r.Cfg.AccessKeyID, r.Cfg.SecretAccessKey, ""),
		),
	)
	if err != nil {
		return nil, err
	}
	r.Client = route53.NewFromConfig(awsCfg)
	return r.Client, nil
}

// ProviderName returns the string "route53" for AWS Route53 providers.
func (r *Route53Provider) ProviderName() string { return "route53" }

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
		HostedZoneId:    aws.String(r.Cfg.HostedZoneID),
		StartRecordName: aws.String(r.Cfg.RecordName),
		StartRecordType: r53types.RRType(r.Cfg.RecordType),
		MaxItems:        aws.Int32(1),
	}
	resp, err := client.ListResourceRecordSets(ctx, input)
	if err != nil {
		return "", err
	}
	for _, record := range resp.ResourceRecordSets {
		if strings.EqualFold(*record.Name, r.Cfg.RecordName+".") && string(record.Type) == r.Cfg.RecordType {
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
		HostedZoneId: aws.String(r.Cfg.HostedZoneID),
		ChangeBatch: &r53types.ChangeBatch{
			Changes: []r53types.Change{
				{
					Action: r53types.ChangeActionUpsert,
					ResourceRecordSet: &r53types.ResourceRecordSet{
						Name:            aws.String(r.Cfg.RecordName),
						Type:            r53types.RRType(r.Cfg.RecordType),
						TTL:             aws.Int64(300),
						ResourceRecords: []r53types.ResourceRecord{{Value: aws.String(ip)}},
					},
				},
			},
		},
	}
	_, err = client.ChangeResourceRecordSets(ctx, input)
	return err
}
