package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds all configuration for JSMon
type Config struct {
	// Telegram configuration
	TelegramToken  string
	TelegramChatID string
	NotifyTelegram bool

	// Slack configuration
	SlackToken     string
	SlackChannelID string
	NotifySlack    bool

	// Discord configuration
	DiscordWebhook string
	NotifyDiscord  bool
}

// Load reads configuration from environment variables
// It loads .env file if present but doesn't error if missing
func Load() (*Config, error) {
	// Try to load .env file (ignore error if it doesn't exist)
	_ = godotenv.Load()

	cfg := &Config{
		TelegramToken:  getEnv("JSMON_TELEGRAM_TOKEN", ""),
		TelegramChatID: getEnv("JSMON_TELEGRAM_CHAT_ID", ""),
		NotifyTelegram: getBoolEnv("JSMON_NOTIFY_TELEGRAM", false),

		SlackToken:     getEnv("JSMON_SLACK_TOKEN", ""),
		SlackChannelID: getEnv("JSMON_SLACK_CHANNEL_ID", ""),
		NotifySlack:    getBoolEnv("JSMON_NOTIFY_SLACK", false),

		DiscordWebhook: getEnv("JSMON_DISCORD_WEBHOOK", ""),
		NotifyDiscord:  getBoolEnv("JSMON_NOTIFY_DISCORD", false),
	}

	// Validate that at least one notification method is configured
	if !cfg.NotifyTelegram && !cfg.NotifySlack && !cfg.NotifyDiscord {
		return nil, fmt.Errorf("at least one notification method must be enabled")
	}

	// Validate Telegram config if enabled
	if cfg.NotifyTelegram {
		if cfg.TelegramToken == "" || cfg.TelegramChatID == "" {
			return nil, fmt.Errorf("telegram enabled but JSMON_TELEGRAM_TOKEN or JSMON_TELEGRAM_CHAT_ID not set")
		}
	}

	// Validate Slack config if enabled
	if cfg.NotifySlack {
		if cfg.SlackToken == "" || cfg.SlackChannelID == "" {
			return nil, fmt.Errorf("slack enabled but JSMON_SLACK_TOKEN or JSMON_SLACK_CHANNEL_ID not set")
		}
	}

	// Validate Discord config if enabled
	if cfg.NotifyDiscord {
		if cfg.DiscordWebhook == "" {
			return nil, fmt.Errorf("discord enabled but JSMON_DISCORD_WEBHOOK not set")
		}
	}

	return cfg, nil
}

// getEnv retrieves an environment variable or returns a default value
func getEnv(key, defaultVal string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultVal
}

// getBoolEnv retrieves a boolean environment variable or returns a default value
func getBoolEnv(key string, defaultVal bool) bool {
	valStr := os.Getenv(key)
	if valStr == "" {
		return defaultVal
	}
	val, err := strconv.ParseBool(valStr)
	if err != nil {
		return defaultVal
	}
	return val
}
