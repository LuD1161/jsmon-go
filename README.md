# JSMon-Go

**JSMon-Go** - JavaScript Change Monitor for Bug Bounty (Go Edition)

A high-performance Go rewrite of the original [JSMon](https://github.com/robre/jsmon) Python tool. Monitor JavaScript files on websites and get notified when they change, helping you discover new attack surfaces and vulnerabilities during bug bounty hunting.

## Features

- 🚀 **Single Binary** - No Python runtime or dependencies needed
- ⚡ **Concurrent Fetching** - Fetch multiple endpoints in parallel using goroutines
- 🔄 **Retry Logic** - Automatic retry with exponential backoff for failed requests
- 📊 **Beautiful HTML Diffs** - Beautified JavaScript with side-by-side comparison
- 🔔 **Multi-Channel Notifications** - Telegram, Slack, and Discord support
- 🐛 **Bug Fixes** - Resolves Discord webhook bug from original Python version
- 🏗️ **Better Architecture** - Clean, modular codebase with proper error handling
- 📦 **Easy Distribution** - Cross-compile for Linux, macOS, and Windows

## Why Go?

### Performance Benefits
- **Faster Execution**: Compiled binary vs. interpreted Python
- **Lower Memory**: ~10MB footprint vs. Python interpreter overhead
- **Concurrent**: Native goroutines for parallel endpoint fetching
- **Single Binary**: No dependency management or virtual environments

### Deployment Benefits
- **No Runtime**: Works on any system without Go/Python installed
- **Cross-Platform**: Easy cross-compilation for any OS/architecture
- **Cron-Friendly**: Single binary, no shell wrapper needed
- **Docker-Ready**: Smaller images, faster startup

## Installation

### Option 1: Download Pre-built Binary

Download the latest release for your platform from the [Releases](https://github.com/aseemshrey/jsmon-go/releases) page:

```bash
# Linux (amd64)
wget https://github.com/aseemshrey/jsmon-go/releases/latest/download/jsmon-linux-amd64
chmod +x jsmon-linux-amd64
mv jsmon-linux-amd64 /usr/local/bin/jsmon

# macOS (Apple Silicon)
curl -L https://github.com/aseemshrey/jsmon-go/releases/latest/download/jsmon-darwin-arm64 -o jsmon
chmod +x jsmon
mv jsmon /usr/local/bin/jsmon

# Windows
# Download jsmon-windows-amd64.exe and add to PATH
```

### Option 2: Build from Source

```bash
# Clone the repository
git clone https://github.com/aseemshrey/jsmon-go.git
cd jsmon-go

# Build
make build

# Or install to GOPATH/bin
make install

# Or build for all platforms
make build-all
```

## Configuration

Create a `.env` file in the directory where you'll run jsmon:

```bash
# Telegram (recommended)
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

**Telegram** (Easiest):
1. Message [@BotFather](https://t.me/botfather) and create a bot with `/newbot`
2. Copy the token
3. Message your bot and visit `https://api.telegram.org/bot<TOKEN>/getUpdates` to get your chat_id

**Slack**:
1. Create a Slack App with `files:write` permission
2. Install to workspace and copy OAuth token
3. Get channel ID from channel details

**Discord**:
1. Server Settings → Integrations → Webhooks
2. Create webhook and copy URL

## Usage

### Adding Targets

Create files in the `targets/` directory with one URL per line:

```bash
# Create targets directory
mkdir targets

# Add endpoints to monitor
cat > targets/example-site << EOF
https://example.com/static/js/app.js
https://example.com/static/js/bundle.js
https://cdn.example.com/scripts/main.js
EOF

# You can organize by program/site
cat > targets/hackerone-target << EOF
https://target.example.com/assets/application.js
https://target.example.com/webpack/vendor.js
EOF
```

Lines starting with `#` are treated as comments and ignored.

### Running JSMon

```bash
# Run once
jsmon

# Or use the binary directly
./bin/jsmon
```

**First run**: All endpoints are enrolled without sending notifications.

**Subsequent runs**: Changes are detected and notifications sent with HTML diff attachments.

### Automated Monitoring with Cron

```bash
# Edit crontab
crontab -e

# Run daily at midnight
0 0 * * * cd /path/to/jsmon-go && ./bin/jsmon >> /var/log/jsmon.log 2>&1

# Run every 6 hours
0 */6 * * * cd /path/to/jsmon-go && ./bin/jsmon >> /var/log/jsmon.log 2>&1

# Run hourly for critical targets
0 * * * * cd /path/to/jsmon-go && ./bin/jsmon >> /var/log/jsmon.log 2>&1
```

**Note**: Unlike the Python version, no shell wrapper is needed - the binary includes everything.

## How It Works

1. **Load Targets**: Reads all files in `targets/` directory
2. **Fetch Content**: Concurrently fetches all endpoints with retry logic
3. **Hash & Compare**: MD5 hash comparison against `jsmon.json` state
4. **Detect Changes**: Identifies new, unchanged, and modified files
5. **Generate Diff**: Creates beautified HTML diff for changed files
6. **Notify**: Sends notifications via configured channels with diff attachment

## Development

```bash
# Install dependencies
make deps

# Build
make build

# Run tests
make test

# Build for all platforms
make build-all

# Clean build artifacts
make clean

# Show all targets
make help
```

## Improvements Over Python Version

### Fixed Bugs
- ✅ Discord webhook actually uses the webhook URL (was hardcoded string)
- ✅ Proper error handling throughout
- ✅ Validates configuration before running

### New Features
- ✅ Concurrent endpoint fetching (much faster for many targets)
- ✅ Retry logic with exponential backoff
- ✅ Better progress output with status symbols
- ✅ Single binary distribution
- ✅ Cross-platform builds out of the box

### Better Architecture
- ✅ Clean separation of concerns (config, storage, fetcher, differ, notifier)
- ✅ Proper error propagation
- ✅ Type safety
- ✅ Testable components

## Project Structure

```
jsmon-go/
├── cmd/
│   └── jsmon/
│       └── main.go           # Application entry point
├── internal/
│   ├── config/              # Configuration management
│   ├── storage/             # State and file storage
│   ├── fetcher/             # HTTP endpoint fetching
│   ├── differ/              # Diff generation
│   └── notifier/            # Notification backends
│       ├── telegram.go
│       ├── slack.go
│       └── discord.go
├── Makefile                 # Build automation
├── go.mod                   # Go dependencies
└── README.md
```

## Contributing

Contributions welcome! Please feel free to submit a Pull Request.

## Original Python Version

This is a Go rewrite of the original Python JSMon by [@r0bre](https://github.com/robre/jsmon).

### Original Contributors
- [@r0bre](https://twitter.com/r0bre) - Core Python version
- [@Yassineaboukir](https://twitter.com/Yassineaboukir) - Slack notifications
- [@seczq](https://twitter.com/seczq) - Discord notifications

## License

MIT License - See LICENSE file for details

## Security Note

This tool is designed for authorized security testing and bug bounty hunting. Ensure you have permission to monitor the targets you configure. Respect rate limits and robots.txt.
