# dynago Documentation

[![Go Reference](https://pkg.go.dev/badge/github.com/aaronlmathis/dynago.svg)](https://pkg.go.dev/github.com/aaronlmathis/dynago)
[![GPLv3 License](https://img.shields.io/badge/License-GPLv3-blue.svg)](https://www.gnu.org/licenses/gpl-3.0)
[![Static Analysis: A+](https://img.shields.io/badge/Static%20Analysis-A%2B-brightgreen)](https://goreportcard.com/report/github.com/aaronlmathis/dynago)

---

# dynago

**dynago** is a cross-platform, resource-light dynamic DNS updater written in Go. It automatically updates DNS records with your current public IP address, ensuring your domain always points to your home, server, or cloud instance—even when your IP changes.

## Features

- **Supports multiple DNS providers:**
  - Cloudflare (with support for the "proxied" flag)
  - AWS Route53
- **Efficient:** Only updates DNS records if your public IP has changed.
- **Configurable:** YAML-based configuration for update interval, IP source, logging, and provider-specific options.
- **Robust logging:** Pretty console output and file logging with log levels.
- **Production-ready:** Systemd service and Makefile for easy deployment.
- **Tested:** All modules have unit tests and detailed GoDoc comments.

## Providers Supported

- **Cloudflare**
  - Supports A/AAAA records
  - Uses API token authentication
  - Respects the `proxied` flag (orange cloud)
- **AWS Route53**
  - Supports A/AAAA records
  - Uses static credentials (access key/secret)

## How It Works

1. Periodically fetches your current public IP from a configurable source (e.g., https://api.ipify.org).
2. Checks the configured DNS record(s) at your provider(s).
3. If the DNS record does not match your current IP, updates it via the provider's API.
4. Logs all actions and errors.

## Installation

### Linux

1. **Build:**
   ```bash
   make
   ```
   This builds the dynago binary to `bin/dynago`.

2. **Configure:**
   ```bash
   sudo mkdir -p /etc/dynago
   sudo cp configs/dynago.yml /etc/dynago/dynago.yml
   ```
   Edit `/etc/dynago/dynago.yml` with your provider credentials and desired settings.

3. **Install:**
   ```bash
   sudo make install
   ```
   This will:
   - Install the binary to `/usr/local/bin/dynago`
   - Install the config to `/etc/dynago/dynago.yml`
   - Install the systemd service to `/etc/systemd/system/dynago.service`
   - Enable the service (but not start it)

4. **Start the Service:**
   ```bash
   sudo systemctl daemon-reload
   sudo systemctl start dynago
   sudo systemctl status dynago
   ```

5. **Logs:**
   ```bash
   sudo journalctl -u dynago -f
   ```

---

### Windows

1. **Build:**
   ```powershell
   go build -o dynago.exe ./cmd/dynago
   ```

2. **Configure:**
   ```powershell
   mkdir C:\ProgramData\dynago
   copy configs\dynago.yml C:\ProgramData\dynago\dynago.yml
   ```
   Edit `C:\ProgramData\dynago\dynago.yml` with your provider credentials and desired settings.

3. **Run manually:**
   ```powershell
   dynago.exe -config C:\ProgramData\dynago\dynago.yml
   ```

4. **(Optional) Install as a service:**
   - Download and install [NSSM](https://nssm.cc/).
   - Run:
     ```powershell
     nssm install dynago "C:\path\to\dynago.exe" -config "C:\ProgramData\dynago\dynago.yml"
     nssm start dynago
     ```

5. **Logs:**
   - Use the `-log` flag to specify a log file, e.g. `-log C:\ProgramData\dynago\dynago.log`.

---

## Configuration

Edit `/etc/dynago/dynago.yml` to set your update interval, IP source, log level, and provider credentials. Example:

```yaml
interval: 5m
ip_source: "https://api.ipify.org"
log_level: "info"
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

- Set `enabled: true` for the provider(s) you want to use.
- For Cloudflare, set `proxied: true` to enable the orange cloud (proxy).

## Advanced

- **Run manually:**
  ```
  ./bin/dynago -config=configs/dynago.yml
  ```
- **Change log file:**
  ```
  ./bin/dynago -config=configs/dynago.yml -log=/var/log/dynago.log
  ```
- **Test:**
  ```
  go test ./...
  ```

## Security
- Store your API tokens and credentials securely.
- The config file should be readable only by the user running dynago.

## Contributing
Pull requests and issues are welcome! Please add tests and GoDoc comments for new features.

## License

This project is licensed under the GNU General Public License v3.0. See [LICENSE](LICENSE) for details.

---

# Go API Documentation

## Internal Packages

See [internal.md](internal.md) for detailed documentation of all internal packages (config, logger, provider, service, utils).

## Command Package

See [cmd.md](cmd.md) for documentation of the main command-line entry point.
