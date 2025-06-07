# dynago Providers System

This document describes the provider system in dynago, how to configure providers, and how to add new providers.

---

## Overview

dynago supports multiple DNS providers (e.g., Cloudflare, AWS Route53) for dynamic DNS updates. The provider system is designed to be extensible: each provider is implemented as a separate Go package under `providers/`, and new providers can be added easily by contributors.

- The provider interface and registry are public and live in `providers/provider.go`.
- Each provider (e.g., Cloudflare, Route53) is implemented in its own subpackage (e.g., `providers/cloudflare/`).
- Each provider defines and unmarshals its own config struct.
- Provider configs are loaded as generic maps and passed to the provider's `New` function for parsing.

---

## Provider Interface

The core interface for all providers is defined in `providers/provider.go`:

```go
// DNSProvider is the interface all DNS providers must implement.
type DNSProvider interface {
    Name() string
    Enabled() bool
    Update(recordName, recordType, ip string) error
}
```

Providers must also register themselves using the registry in `provider.go`:

```go
// RegisterProvider registers a provider constructor by name.
func RegisterProvider(name string, constructor ProviderConstructor)
```

---

## Adding a New Provider

1. **Create a new package:**
   - Add a new folder under `providers/` (e.g., `providers/example/`).
   - Implement a `Config` struct for your provider's settings.
   - Implement the `DNSProvider` interface.
   - Register your provider in an `init()` function.

2. **Example skeleton:**

```go
// providers/example/example.go
package example

type Config struct {
    Enabled bool   `yaml:"enabled"`
    ApiKey  string `yaml:"api_key"`
    // ...other fields...
}

type ExampleProvider struct {
    cfg Config
}

func (p *ExampleProvider) Name() string    { return "example" }
func (p *ExampleProvider) Enabled() bool   { return p.cfg.Enabled }
func (p *ExampleProvider) Update(recordName, recordType, ip string) error {
    // ...implementation...
    return nil
}

func New(config map[string]any) (providers.DNSProvider, error) {
    var cfg Config
    // Unmarshal config map to struct (see ConfigFromMap helper)
    if err := providers.ConfigFromMap(config, &cfg); err != nil {
        return nil, err
    }
    return &ExampleProvider{cfg: cfg}, nil
}

func init() {
    providers.RegisterProvider("example", New)
}
```

3. **Write tests:**
   - Add a `example_test.go` file with unit tests for your provider.

---

## Provider Configuration in YAML

Provider configs are specified under the `providers:` key in your YAML config file. Each provider's config is a map of options, passed directly to the provider's `New` function.

Example:

```yaml
providers:
  cloudflare:
    enabled: true
    api_token: "your-cloudflare-api-token"
    zone_id: "your-zone-id"
    record_name: "home.example.com"
    record_type: "A"
    proxied: true
  route53:
    enabled: false
    access_key_id: "AWS_ACCESS_KEY_ID"
    secret_access_key: "AWS_SECRET_ACCESS_KEY"
    hosted_zone_id: "Z1D633PJN98FT9"
    record_name: "home.example.com"
    record_type: "A"
    region: "us-east-1"
```

- Each provider can define its own config fields.
- Only providers with `enabled: true` will be used.

---

## Provider Registry and Discovery

- The registry in `provider.go` allows dynago to discover and instantiate all registered providers at runtime.
- Providers are enabled/disabled via config.
- The registry makes it easy to add new providers without modifying core code.

---

## See Also
- Example provider implementations: `providers/cloudflare/`, `providers/route53/`
- Main configuration documentation: [README.md](README.md)
