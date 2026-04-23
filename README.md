# :satellite: Server Monitor

A lightweight Go CLI tool that collects OS-level health metrics from a Linux server and sends an HTML email alert when any metric exceeds its threshold. Designed to run periodically via systemd timer (or cron).

## Monitored Metrics

| Metric | Default Threshold |
|---|---|
| CPU usage | 75% |
| Memory usage | 75% |
| Board temperature | 60°C |
| Load average (5min / 15min) | 75% |

Additional stats included in alerts: network Rx/Tx, uptime, hostname.

## Requirements

- **Go** 1.26+ (for building)
- **Linux** (relies on `/proc` and `/sys` for stats)
- A local MTA (`sendmail`) configured for email delivery
- **golangci-lint** and **staticcheck** (for linting)

## Configuration

All configuration is done via environment variables, typically provided through a `.env` file loaded by systemd's `EnvironmentFile=` directive.

```env
# Required - email
SOURCE_EMAIL_ADDRESS=alerts@example.com
TARGET_EMAIL_ADDRESS=admin@example.com

# Optional - thresholds
SERVER_NAME=my-server
THRESHOLD_TEMP=60
THRESHOLD_MEM=75
THRESHOLD_CPU=75
```

## Tooling

- [Go](https://go.dev/dl/) 1.26+
- [GNU Make](https://www.gnu.org/software/make/)

```bash
go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest
go install honnef.co/go/tools/cmd/staticcheck@latest
```

## Build

```bash
# Build for all targets (linux/amd64 + linux/arm64)
make build

# Build for a specific architecture
make build-amd64
make build-arm64

# Lint, test, and build
make
```

Binaries are output to `dist/`.

## Deployment

1. Copy the binary to your server:
   ```bash
   scp dist/server-monitor-linux-arm64 user@server:/home/Server-Monitor/server-monitor
   ```

2. Create your `.env` file on the server:
   ```bash
   vi /home/Server-Monitor/.env
   ```

3. Install the systemd units:
   ```bash
   sudo cp systemd/server-monitor.service /etc/systemd/system/
   sudo cp systemd/server-monitor.timer /etc/systemd/system/
   sudo systemctl daemon-reload
   sudo systemctl enable --now server-monitor.timer
   ```

4. Verify:
   ```bash
   # Check timer status
   systemctl status server-monitor.timer

   # Run manually
   sudo systemctl start server-monitor.service

   # View logs
   journalctl -u server-monitor.service
   ```

## Project Structure

```
server-monitor/
├── main.go                        # Entry point, config, alerting logic
├── internal/monitor/
│   ├── stats.go                   # OS metrics collection + utilities
│   ├── mail.go                    # Email alert delivery via sendmail
│   └── mail_template.go           # HTML email template
├── systemd/
│   ├── server-monitor.service     # Oneshot service unit
│   └── server-monitor.timer       # Hourly timer
├── Makefile                       # Build, lint, test targets
└── dist/                          # Build output (gitignored)
```
