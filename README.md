<div align="center">

# JSMon-Go

**JavaScript Change Monitor for Bug Bounty Hunting**

[![Go Version](https://img.shields.io/github/go-mod/go-version/LuD1161/jsmon-go)](https://go.dev/)
[![Release](https://img.shields.io/github/v/release/LuD1161/jsmon-go)](https://github.com/LuD1161/jsmon-go/releases)
[![License](https://img.shields.io/github/license/LuD1161/jsmon-go?)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/LuD1161/jsmon-go?)](https://goreportcard.com/report/github.com/LuD1161/jsmon-go)

A high-performance Go rewrite of the original [JSMon](https://github.com/robre/jsmon) Python tool. Monitor JavaScript files on websites and get notified when they change, helping you discover new attack surfaces and vulnerabilities during bug bounty hunting.

[Features](#features) • [Installation](#installation) • [Quick Start](#quick-start) • [Documentation](#documentation) • [Docker](#docker)

</div>

---

## Features

- 🚀 **Single Binary** - No Python runtime or dependencies needed
- ⚡ **10x Faster** - Concurrent fetching with goroutines vs sequential Python
- 🔄 **Retry Logic** - Automatic retry with exponential backoff for failed requests
- 📊 **Beautiful HTML Diffs** - Beautified JavaScript with side-by-side comparison
- 🔔 **Multi-Channel Notifications** - Telegram, Slack, and Discord support
- 🐛 **Bug Fixes** - Resolves Discord webhook bug from original Python version
- 🏗️ **Better Architecture** - Clean, modular codebase with proper error handling
- 📦 **Easy Distribution** - Cross-compile for Linux, macOS, Windows, and Docker
- 🔒 **Minimal Attack Surface** - Distroless Docker images for security

## Why Go?

| Feature | Python Version | Go Version |
|---------|---------------|------------|
| **Execution** | Interpreted (~500ms startup) | Compiled (~10ms startup) |
| **Concurrency** | Sequential | Parallel (goroutines) |
| **Binary Size** | N/A (requires runtime) | 8.9MB (includes everything) |
| **Memory** | ~50-100MB | ~10-20MB |
| **Speed (100 endpoints)** | ~100s | ~10s (10x faster) |
| **Dependencies** | pip packages + venv | Embedded in binary |
| **Discord Bug** | Present | Fixed ✅ |

## Quick Start

```bash
# Download and extract binary (Linux amd64 example)
VERSION=1.0.1  # Check https://github.com/LuD1161/jsmon-go/releases for latest
curl -L https://github.com/LuD1161/jsmon-go/releases/download/v${VERSION}/jsmon-go_${VERSION}_Linux_x86_64.tar.gz | tar xz
sudo mv jsmon /usr/local/bin/jsmon

# Create configuration
cat > .env << EOF
JSMON_NOTIFY_TELEGRAM=true
JSMON_TELEGRAM_TOKEN=your_bot_token
JSMON_TELEGRAM_CHAT_ID=your_chat_id
EOF

# Add targets to monitor
mkdir targets
cat > targets/my-targets << EOF
https://example.com/static/js/app.js
https://example.com/static/js/bundle.js
EOF

# Run (first run enrolls endpoints without notifications)
jsmon

# Run again to detect changes
jsmon
```

## Installation

### Option 1: Docker (Recommended for Production)

```bash
# Pull the image
docker pull ghcr.io/LuD1161/jsmon-go:latest

# Run with mounted volumes
docker run --rm \
  -v $(pwd)/targets:/app/targets \
  -v $(pwd)/downloads:/app/downloads \
  -v $(pwd)/jsmon.json:/app/jsmon.json \
  -v $(pwd)/.env:/app/.env \
  ghcr.io/LuD1161/jsmon-go:latest
```

### Option 2: Pre-built Binary

Download the latest release for your platform from the [Releases](https://github.com/LuD1161/jsmon-go/releases) page.

**Note**: Replace `VERSION` with the latest version (e.g., `1.0.1`) from the [releases page](https://github.com/LuD1161/jsmon-go/releases).

**Linux:**
```bash
# amd64
VERSION=1.0.1
curl -L https://github.com/LuD1161/jsmon-go/releases/download/v${VERSION}/jsmon-go_${VERSION}_Linux_x86_64.tar.gz | tar xz
sudo mv jsmon /usr/local/bin/jsmon

# arm64
VERSION=1.0.1
curl -L https://github.com/LuD1161/jsmon-go/releases/download/v${VERSION}/jsmon-go_${VERSION}_Linux_arm64.tar.gz | tar xz
sudo mv jsmon /usr/local/bin/jsmon
```

**macOS:**
```bash
# Apple Silicon (M1/M2/M3)
VERSION=1.0.1
curl -L https://github.com/LuD1161/jsmon-go/releases/download/v${VERSION}/jsmon-go_${VERSION}_Darwin_arm64.tar.gz | tar xz
sudo mv jsmon /usr/local/bin/jsmon

# Intel
VERSION=1.0.1
curl -L https://github.com/LuD1161/jsmon-go/releases/download/v${VERSION}/jsmon-go_${VERSION}_Darwin_x86_64.tar.gz | tar xz
sudo mv jsmon /usr/local/bin/jsmon
```

**Windows:**
```powershell
# Download and extract (replace 1.0.1 with latest version)
$VERSION = "1.0.1"
curl -L -o jsmon-go.zip "https://github.com/LuD1161/jsmon-go/releases/download/v$VERSION/jsmon-go_${VERSION}_Windows_x86_64.zip"
Expand-Archive jsmon-go.zip -DestinationPath .
# Move jsmon.exe to a directory in your PATH
```

### Option 3: Build from Source

```bash
# Clone the repository
git clone https://github.com/LuD1161/jsmon-go.git
cd jsmon-go

# Build for current platform
make build

# Or install to GOPATH/bin
make install

# Or build for all platforms
make build-all
```

### Option 4: Go Install

```bash
go install github.com/LuD1161/jsmon-go/cmd/jsmon@latest
```

## Configuration

Create a `.env` file in your working directory:

```bash
# Telegram (Recommended - easiest to set up)
JSMON_NOTIFY_TELEGRAM=true
JSMON_TELEGRAM_TOKEN=your_bot_token_here
JSMON_TELEGRAM_CHAT_ID=your_chat_id_here

# Slack
JSMON_NOTIFY_SLACK=false
JSMON_SLACK_TOKEN=xoxb-your-token
JSMON_SLACK_CHANNEL_ID=C01234567

# Discord
JSMON_NOTIFY_DISCORD=false
JSMON_DISCORD_WEBHOOK=https://discord.com/api/webhooks/...
```

### Getting Notification Credentials

<details>
<summary><b>Telegram</b> (Easiest - Recommended)</summary>

1. Message [@BotFather](https://t.me/botfather) and create a bot with `/newbot`
2. Copy the token provided
3. Message your bot with any text
4. Visit `https://api.telegram.org/bot<YOUR_TOKEN>/getUpdates` to get your `chat_id`
5. Add both values to your `.env` file

</details>

<details>
<summary><b>Slack</b></summary>

1. Create a Slack App at https://api.slack.com/apps
2. Add `files:write` and `chat:write` permissions
3. Install to your workspace
4. Copy the OAuth token (starts with `xoxb-`)
5. Get your channel ID from channel details (right-click channel → View channel details)

</details>

<details>
<summary><b>Discord</b></summary>

1. Go to Server Settings → Integrations → Webhooks
2. Create a new webhook
3. Copy the webhook URL
4. Add to your `.env` file

</details>

## Usage

### Adding Targets

Create files in the `targets/` directory with one URL per line:

```bash
mkdir targets

# Add endpoints to monitor
cat > targets/example-site << EOF
# Main application bundle
https://example.com/static/js/app.js
https://example.com/static/js/bundle.js

# CDN resources
https://cdn.example.com/scripts/main.js
EOF

# Organize by program/site
cat > targets/hackerone-program << EOF
https://target.example.com/assets/application.js
https://target.example.com/webpack/vendor.js
EOF
```

**Notes:**
- Lines starting with `#` are treated as comments
- Empty lines are ignored
- You can organize targets across multiple files

### Running JSMon

```bash
# Run once
jsmon

# With custom working directory
cd /path/to/monitoring && jsmon
```

**First run**: All endpoints are enrolled without sending notifications.

**Subsequent runs**: Changes are detected and notifications sent with beautified HTML diff attachments.

### Example Output

```
JSMon - Web File Monitor (Go Edition)
Monitoring 5 endpoint(s)...
✓ Telegram notifications enabled
  No change: https://example.com/app.js
⚠ Change detected: https://example.com/bundle.js (abc123def4 → xyz789ghi0)
  ✓ Notification sent
⊕ New endpoint enrolled: https://example.com/new.js
✗ Failed to fetch https://down.example.com/app.js: connection timeout

==================================================
Summary: 5 endpoints monitored
  • 1 new endpoint(s) enrolled
  • 1 change(s) detected and notified
  • 1 error(s) occurred
```

### Automated Monitoring

**Cron (Linux/macOS):**
```bash
# Edit crontab
crontab -e

# Run every hour
0 * * * * cd /path/to/jsmon && /usr/local/bin/jsmon >> /var/log/jsmon.log 2>&1

# Run every 6 hours
0 */6 * * * cd /path/to/jsmon && /usr/local/bin/jsmon >> /var/log/jsmon.log 2>&1
```

**Systemd Timer (Linux):**
```ini
# /etc/systemd/system/jsmon.timer
[Unit]
Description=JSMon Change Monitor

[Timer]
OnCalendar=hourly
Persistent=true

[Install]
WantedBy=timers.target
```

**Docker Compose + Cron:**
```yaml
version: '3.8'
services:
  jsmon:
    image: ghcr.io/LuD1161/jsmon-go:latest
    volumes:
      - ./targets:/app/targets
      - ./downloads:/app/downloads
      - ./jsmon.json:/app/jsmon.json
      - ./.env:/app/.env
    restart: "no"
```

Then run via cron:
```bash
0 * * * * cd /path/to/jsmon && docker-compose run --rm jsmon
```

## Docker

### Using Pre-built Images

```bash
# Pull the latest image
docker pull ghcr.io/LuD1161/jsmon-go:latest

# Or a specific version
docker pull ghcr.io/LuD1161/jsmon-go:v1.0.1

# Run with mounted volumes
docker run --rm \
  -v $(pwd)/targets:/app/targets \
  -v $(pwd)/downloads:/app/downloads \
  -v $(pwd)/jsmon.json:/app/jsmon.json \
  -v $(pwd)/.env:/app/.env \
  ghcr.io/LuD1161/jsmon-go:latest
```

### Building Locally

```bash
# Build the image from source (uses Dockerfile.standalone)
docker build -f Dockerfile.standalone -t jsmon-go .

# Run it
docker run --rm \
  -v $(pwd)/targets:/app/targets \
  -v $(pwd)/downloads:/app/downloads \
  -v $(pwd)/jsmon.json:/app/jsmon.json \
  -v $(pwd)/.env:/app/.env \
  jsmon-go
```

**Note**: The default `Dockerfile` is optimized for GoReleaser (expects pre-built binary). For building from source, use `Dockerfile.standalone`.

### Why Distroless?

The Docker image uses Google's `distroless` base image:
- **Security**: No shell, no package manager - minimal attack surface
- **Size**: ~20MB total (vs ~100MB+ with Alpine)
- **Performance**: Only includes runtime dependencies

## How It Works

```
┌─────────────┐
│  targets/   │  Load endpoint URLs
└──────┬──────┘
       │
       ▼
┌─────────────────────────────┐
│  Concurrent HTTP Fetching   │  Goroutines + retry logic
│  (with exponential backoff) │
└──────────┬──────────────────┘
           │
           ▼
    ┌──────────────┐
    │  MD5 Hashing │  10-char truncated hash
    └──────┬───────┘
           │
           ▼
    ┌─────────────────┐
    │ Compare to State│  Check jsmon.json
    │  (jsmon.json)   │
    └────────┬────────┘
             │
        ┌────┴─────┐
        │          │
        ▼          ▼
    No Change   Changed!
        │          │
        │          ▼
        │    ┌──────────────┐
        │    │ Beautify JS  │
        │    │ Generate Diff│
        │    └──────┬───────┘
        │           │
        │           ▼
        │    ┌─────────────────────┐
        │    │ Send Notifications  │
        │    │ (Telegram/Slack/    │
        │    │  Discord)           │
        │    └─────────────────────┘
        │
        └──────────┬──────────────┘
                   │
                   ▼
            ┌──────────────┐
            │  Save State  │
            │ & Content    │
            └──────────────┘
```

## Development

```bash
# Install dependencies
make deps

# Build for current platform
make build

# Run tests
make test

# Build for all platforms (Linux, macOS, Windows)
make build-all

# Clean build artifacts
make clean

# Show all available targets
make help
```

### Project Structure

```
jsmon-go/
├── cmd/
│   └── jsmon/
│       └── main.go              # Application entry point
├── internal/                     # Private packages
│   ├── config/                  # Environment configuration
│   │   └── config.go            # .env loading & validation
│   ├── storage/                 # State management
│   │   └── storage.go           # jsmon.json & downloads/
│   ├── fetcher/                 # HTTP operations
│   │   └── fetcher.go           # Concurrent fetching + retry
│   ├── differ/                  # Diff generation
│   │   └── differ.go            # JS beautification + HTML diff
│   └── notifier/                # Notifications
│       ├── notifier.go          # Interface
│       ├── telegram.go          # Telegram API
│       ├── slack.go             # Slack API
│       └── discord.go           # Discord webhooks
├── Dockerfile                    # Distroless container image
├── .goreleaser.yml              # Release automation
├── Makefile                     # Build automation
├── go.mod                       # Go dependencies
└── README.md
```

## Performance Comparison

**Scenario**: Monitoring 100 JavaScript endpoints

| Metric | Python | Go | Improvement |
|--------|--------|----|-----------:|
| Execution Time | ~100s | ~10s | **10x faster** |
| Memory Usage | ~80MB | ~15MB | **5x less** |
| Binary Size | N/A | 8.9MB | Portable |
| Cold Start | ~500ms | ~10ms | **50x faster** |
| Concurrent Requests | No | Yes | ✅ |
| Retry Logic | No | Yes | ✅ |

## Improvements Over Python Version

### Fixed Bugs
- ✅ **Discord webhook** actually uses the webhook URL (was hardcoded string `"DISCORD_WEBHOOK_URL"`)
- ✅ **Proper error handling** - failed endpoints don't stop entire run
- ✅ **Configuration validation** before execution

### New Features
- ✅ **Concurrent fetching** with goroutines (10x faster for many targets)
- ✅ **Retry logic** with exponential backoff (3 retries, 30s timeout)
- ✅ **Better progress output** with status symbols (✓, ✗, ⊕, ⚠)
- ✅ **Single binary distribution** - no runtime dependencies
- ✅ **Cross-platform builds** - Linux, macOS, Windows, Docker
- ✅ **Docker support** with minimal distroless images

### Better Architecture
- ✅ **Clean separation of concerns** - config, storage, fetcher, differ, notifier
- ✅ **Interface-based design** - easy to test and extend
- ✅ **Type safety** - compile-time error checking
- ✅ **Proper error propagation** throughout the stack

## Contributing

Contributions are welcome! Here's how you can help:

### Reporting Bugs
1. Check existing [issues](https://github.com/LuD1161/jsmon-go/issues)
2. Create a new issue with detailed reproduction steps
3. Include your OS, Go version, and JSMon version

### Suggesting Features
1. Open an issue with the `enhancement` label
2. Describe the use case and expected behavior
3. Be open to discussion and feedback

### Submitting Pull Requests
1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Run tests: `make test`
5. Commit with clear messages
6. Push to your fork
7. Open a Pull Request with a detailed description

### Development Setup
```bash
git clone https://github.com/LuD1161/jsmon-go.git
cd jsmon-go
make deps
make build
```

## Credits

This is a Go rewrite of the original Python [JSMon](https://github.com/robre/jsmon) by [@r0bre](https://github.com/robre).

### Original Contributors
- [@r0bre](https://twitter.com/r0bre) - Core Python version
- [@Yassineaboukir](https://twitter.com/Yassineaboukir) - Slack notifications
- [@seczq](https://twitter.com/seczq) - Discord notifications

## License

MIT License - See [LICENSE](LICENSE) file for details.

## Security & Responsible Use

⚠️ **Important**: This tool is designed for **authorized security testing and bug bounty hunting only**.

- Ensure you have permission to monitor the targets you configure
- Respect rate limits and avoid overwhelming servers
- Follow responsible disclosure practices
- Comply with bug bounty program rules
- Don't use for unauthorized access or malicious purposes

## Support

- 📖 [Documentation](https://github.com/LuD1161/jsmon-go#readme)
- 🐛 [Report Issues](https://github.com/LuD1161/jsmon-go/issues)
- 💬 [Discussions](https://github.com/LuD1161/jsmon-go/discussions)

---

<div align="center">

**Built with ❤️ for the bug bounty community**

[⬆ Back to top](#jsmon-go)

</div>
