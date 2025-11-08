package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/LuD1161/jsmon-go/internal/config"
	"github.com/LuD1161/jsmon-go/internal/differ"
	"github.com/LuD1161/jsmon-go/internal/fetcher"
	"github.com/LuD1161/jsmon-go/internal/notifier"
	"github.com/LuD1161/jsmon-go/internal/storage"
)

func main() {
	fmt.Println("JSMon - Web File Monitor (Go Edition)")

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Configuration error: %v", err)
	}

	// Initialize storage
	store := storage.NewDefault()
	if err := store.Initialize(); err != nil {
		log.Fatalf("Storage initialization error: %v", err)
	}

	// Load current state
	state, err := store.LoadState()
	if err != nil {
		log.Fatalf("Failed to load state: %v", err)
	}

	// Load endpoints from targets directory
	endpoints, err := fetcher.LoadEndpoints(fetcher.TargetsDir)
	if err != nil {
		log.Fatalf("Failed to load endpoints: %v", err)
	}

	if len(endpoints) == 0 {
		fmt.Println("No endpoints found in targets/ directory")
		fmt.Println("Add URLs to files in targets/ (one per line) and run again")
		return
	}

	fmt.Printf("Monitoring %d endpoint(s)...\n", len(endpoints))

	// Setup notifiers
	var notifiers []notifier.Notifier

	if cfg.NotifyTelegram {
		chatID, err := strconv.ParseInt(cfg.TelegramChatID, 10, 64)
		if err != nil {
			log.Fatalf("Invalid Telegram chat ID: %v", err)
		}
		tgNotifier, err := notifier.NewTelegramNotifier(cfg.TelegramToken, chatID)
		if err != nil {
			log.Fatalf("Failed to create Telegram notifier: %v", err)
		}
		notifiers = append(notifiers, tgNotifier)
		fmt.Println("✓ Telegram notifications enabled")
	}

	if cfg.NotifySlack {
		slackNotifier := notifier.NewSlackNotifier(cfg.SlackToken, cfg.SlackChannelID)
		notifiers = append(notifiers, slackNotifier)
		fmt.Println("✓ Slack notifications enabled")
	}

	if cfg.NotifyDiscord {
		discordNotifier := notifier.NewDiscordNotifier(cfg.DiscordWebhook)
		notifiers = append(notifiers, discordNotifier)
		fmt.Println("✓ Discord notifications enabled")
	}

	multiNotifier := notifier.NewMultiNotifier(notifiers...)

	// Create fetcher and differ
	f := fetcher.NewDefault()
	d := differ.New()

	// Fetch all endpoints
	results := f.FetchAll(endpoints)

	// Process each result
	changesDetected := 0
	newEndpoints := 0
	errors := 0

	for _, result := range results {
		if result.Error != nil {
			fmt.Printf("✗ Failed to fetch %s: %v\n", result.Endpoint, result.Error)
			errors++
			continue
		}

		// Calculate hash of fetched content
		newHash := storage.HashContent(result.Content)

		// Get previous hash
		prevHash := store.GetPreviousHash(state, result.Endpoint)

		if prevHash == "" {
			// New endpoint - enroll it
			store.AddHash(state, result.Endpoint, newHash)
			if err := store.SaveEndpointContent(newHash, result.Content); err != nil {
				fmt.Printf("✗ Failed to save content for %s: %v\n", result.Endpoint, err)
				errors++
				continue
			}
			fmt.Printf("⊕ New endpoint enrolled: %s\n", result.Endpoint)
			newEndpoints++
			continue
		}

		if prevHash == newHash {
			// No change
			fmt.Printf("  No change: %s\n", result.Endpoint)
			continue
		}

		// Change detected!
		fmt.Printf("⚠ Change detected: %s (%s → %s)\n", result.Endpoint, prevHash, newHash)
		changesDetected++

		// Save new version
		if err := store.SaveEndpointContent(newHash, result.Content); err != nil {
			fmt.Printf("✗ Failed to save new content: %v\n", err)
			errors++
			continue
		}

		// Update state
		store.AddHash(state, result.Endpoint, newHash)

		// Load old content
		oldContent, err := store.LoadEndpointContent(prevHash)
		if err != nil {
			fmt.Printf("✗ Failed to load old content: %v\n", err)
			errors++
			continue
		}

		// Get file sizes
		oldSize, err := store.GetFileSize(prevHash)
		if err != nil {
			fmt.Printf("✗ Failed to get old file size: %v\n", err)
			errors++
			continue
		}

		newSize, err := store.GetFileSize(newHash)
		if err != nil {
			fmt.Printf("✗ Failed to get new file size: %v\n", err)
			errors++
			continue
		}

		// Generate diff
		diffHTML, err := d.GenerateHTMLDiff(oldContent, result.Content)
		if err != nil {
			fmt.Printf("✗ Failed to generate diff: %v\n", err)
			errors++
			continue
		}

		// Send notification
		notification := &notifier.ChangeNotification{
			Endpoint: result.Endpoint,
			OldHash:  prevHash,
			NewHash:  newHash,
			OldSize:  oldSize,
			NewSize:  newSize,
			DiffHTML: diffHTML,
		}

		if err := multiNotifier.Notify(notification); err != nil {
			fmt.Printf("✗ Failed to send notification: %v\n", err)
			errors++
		} else {
			fmt.Printf("  ✓ Notification sent\n")
		}
	}

	// Save updated state
	if err := store.SaveState(state); err != nil {
		log.Fatalf("Failed to save state: %v", err)
	}

	// Summary
	fmt.Println("\n" + strings.Repeat("=", 50))
	fmt.Printf("Summary: %d endpoints monitored\n", len(endpoints))
	if newEndpoints > 0 {
		fmt.Printf("  • %d new endpoint(s) enrolled\n", newEndpoints)
	}
	if changesDetected > 0 {
		fmt.Printf("  • %d change(s) detected and notified\n", changesDetected)
	}
	if errors > 0 {
		fmt.Printf("  • %d error(s) occurred\n", errors)
		os.Exit(1)
	}

	fmt.Println("Done!")
}
