# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

JSMon-Go is a Go rewrite of the Python JSMon tool for monitoring JavaScript file changes during bug bounty hunting. It fetches endpoints, detects changes via MD5 hashing, and sends notifications with beautified HTML diffs via Telegram, Slack, or Discord.

**Key improvements over Python version**: Single binary distribution, concurrent fetching, retry logic, fixed Discord webhook bug, better error handling.

## Commands

### Building
```bash
# Build for current platform
make build          # Output: bin/jsmon

# Build for all platforms (Linux, macOS, Windows)
make build-all      # Output: bin/jsmon-{os}-{arch}

# Install to GOPATH/bin
make install
```

### Development
```bash
# Install dependencies and tidy go.mod
make deps

# Run without building binary
go run cmd/jsmon/main.go

# Run with build
make run

# Run tests
make test

# Clean build artifacts
make clean
```

### Testing Changes
```bash
# Create test target
echo "https://cdnjs.cloudflare.com/ajax/libs/jquery/3.6.0/jquery.min.js" > targets/test

# Create .env if needed
cat > .env << EOF
JSMON_NOTIFY_TELEGRAM=true
JSMON_TELEGRAM_TOKEN=your_token
JSMON_TELEGRAM_CHAT_ID=your_chat_id
EOF

# Run to enroll (no notification)
go run cmd/jsmon/main.go

# Change the URL version and run again to test change detection
```

## Architecture

### Module Structure
Following Go's standard project layout with `cmd/` for entry points and `internal/` for private packages.

### Package Overview

**[internal/config](internal/config/config.go)** - Environment configuration
- Loads `.env` file using `github.com/joho/godotenv`
- Validates that at least one notification method is enabled
- Validates credentials for enabled notification methods
- Type: `Config` struct with all settings

**[internal/storage](internal/storage/storage.go)** - State and file management
- **jsmon.json**: Maps `endpoint URL → []hash` (hash history array)
- **downloads/{hash}**: Raw content files, named by first 10 chars of MD5
- Key methods: `LoadState()`, `SaveState()`, `SaveEndpointContent()`, `LoadEndpointContent()`
- `HashContent()` function: MD5 hash truncated to 10 chars (matches Python version)

**[internal/fetcher](internal/fetcher/fetcher.go)** - HTTP endpoint fetching
- Concurrent fetching via goroutines + WaitGroup
- Retry logic with exponential backoff (default: 3 retries, 30s timeout)
- `LoadEndpoints()`: Reads all files in `targets/`, parses line-separated URLs
- `FetchAll()`: Fetches all endpoints in parallel, returns `[]FetchResult`
- URL validation via regex (same pattern as Python version)

**[internal/differ](internal/differ/differ.go)** - HTML diff generation
- Uses `github.com/ditashi/jsbeautifier-go` for JS beautification
- Uses `github.com/sergi/go-diff` for diff algorithm
- Generates styled HTML with side-by-side comparison
- Falls back to raw content if beautification fails

**[internal/notifier](internal/notifier/)** - Multi-channel notifications
- Interface-based design: all notifiers implement `Notifier` interface
- **telegram.go**: Uses `github.com/go-telegram-bot-api/telegram-bot-api/v5`, sends via `sendDocument` with HTML caption
- **slack.go**: Uses `github.com/slack-go/slack`, uploads file with `files_upload`
- **discord.go**: Standard HTTP webhook POST with multipart form (fixes Python bug)
- `MultiNotifier`: Wraps multiple notifiers, sends to all enabled channels

**[cmd/jsmon/main.go](cmd/jsmon/main.go)** - Application orchestration
1. Load config and validate
2. Initialize storage (create directories/files if needed)
3. Load endpoint list from `targets/`
4. Setup enabled notifiers
5. Fetch all endpoints concurrently
6. Process results: detect new/unchanged/changed
7. For changes: save content, generate diff, send notification
8. Save updated state

### Data Flow
```
targets/*.txt → LoadEndpoints() → FetchAll() [concurrent] → Compare hashes
                                       ↓
                             Changed? → Generate diff → MultiNotifier
                                       ↓
                             Save state & content
```

### Dependencies

**Required third-party packages** (see [go.mod](go.mod)):
- `github.com/joho/godotenv` - .env file loading
- `github.com/go-telegram-bot-api/telegram-bot-api/v5` - Telegram API
- `github.com/slack-go/slack` - Slack API
- `github.com/ditashi/jsbeautifier-go` - JavaScript beautification
- `github.com/sergi/go-diff` - Diff algorithm

Run `make deps` or `go mod tidy` after adding new imports.

### State Management

**jsmon.json format**:
```json
{
  "https://example.com/app.js": ["abc123def4", "xyz789ghi0"],
  "https://cdn.example.com/bundle.js": ["fed456cba9"]
}
```
- Each endpoint has array of hashes (history)
- Latest hash is last item in array
- New endpoints get array with single hash
- Empty `{}` if no endpoints enrolled yet

**downloads/ structure**:
- Files named by hash: `abc123def4`, `xyz789ghi0`, etc.
- Contains raw fetched content (not beautified)
- Used for diff generation and size calculation

## Key Design Decisions

### Why first 10 chars of MD5?
Matches Python version for compatibility. Balance between uniqueness (collision unlikely with small file count) and readable filenames.

### Why goroutines for fetching?
Bug bounty hunters often monitor 100+ endpoints. Sequential fetching (Python) is slow. Concurrent fetching with goroutines provides massive speedup with minimal code complexity.

### Why retry logic?
Networks are unreliable. Temporary failures shouldn't require manual re-runs. Exponential backoff prevents hammering servers.

### Why separate notifier interface?
- Testability: Mock notifiers for testing
- Extensibility: Add new channels without modifying core logic
- Clean separation: Notification logic isolated from detection logic

### Why internal/ not pkg/?
Packages in `internal/` cannot be imported by external projects. This prevents accidental API surface area since JSMon-Go is a CLI tool, not a library.

## Common Tasks

### Adding a new notification channel
1. Create `internal/notifier/newservice.go`
2. Implement `Notifier` interface (`Notify()`, `Name()`)
3. Add config fields to `internal/config/config.go`
4. Add initialization in `cmd/jsmon/main.go` (around lines 62-80)
5. Update README.md with setup instructions

### Changing hash algorithm
- Modify `storage.HashContent()` function
- Keep truncation at 10 chars for filename compatibility
- Consider: Migration path for existing jsmon.json files

### Adding command-line flags
- Import `flag` package in main.go
- Define flags before `config.Load()`
- Consider: Override .env values with flag values if provided

### Debugging notification issues
- Check config validation in `config.Load()`
- Each notifier returns specific errors
- Test individual notifiers by temporarily disabling others in .env

## Differences from Python Version

**Fixed bugs**:
- [Discord webhook bug](https://github.com/robre/jsmon/blob/master/jsmon.py#L158): Was hardcoded string `"DISCORD_WEBHOOK_URL"`, now uses actual variable

**Architectural improvements**:
- Concurrent fetching (Python: sequential)
- Retry logic (Python: none)
- Proper error handling (Python: basic)
- Package organization (Python: single file)
- Type safety (Go: compile-time, Python: runtime)

**Behavioral changes**:
- Better progress output with status symbols (✓, ✗, ⊕, ⚠)
- No shell wrapper needed for cron
- Failed requests don't stop entire run (continues with other endpoints)

## Testing Strategy

Currently no automated tests. When adding tests:

**Unit tests**:
- `storage_test.go`: Test hash calculation, state load/save
- `fetcher_test.go`: Test URL validation, endpoint loading
- `differ_test.go`: Test diff generation with sample JS

**Integration tests**:
- Mock HTTP server for testing fetcher
- Mock notifiers for testing notification flow
- Temp directory for testing storage operations

Use Go's standard `testing` package. Run with `go test ./...` or `make test`.
