package notifier

import (
	"bytes"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// TelegramNotifier sends notifications via Telegram
type TelegramNotifier struct {
	bot    *tgbotapi.BotAPI
	chatID int64
}

// NewTelegramNotifier creates a new Telegram notifier
func NewTelegramNotifier(token string, chatID int64) (*TelegramNotifier, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, fmt.Errorf("failed to create telegram bot: %w", err)
	}

	return &TelegramNotifier{
		bot:    bot,
		chatID: chatID,
	}, nil
}

// Name returns the notifier name
func (t *TelegramNotifier) Name() string {
	return "Telegram"
}

// Notify sends a change notification via Telegram
func (t *TelegramNotifier) Notify(notification *ChangeNotification) error {
	// Create HTML caption with formatting
	caption := fmt.Sprintf(
		"%s has been updated from <code>%s</code> (<b>%d</b> bytes) to <code>%s</code> (<b>%d</b> bytes)",
		notification.Endpoint,
		notification.OldHash,
		notification.OldSize,
		notification.NewHash,
		notification.NewSize,
	)

	// Create document from HTML diff
	diffBytes := bytes.NewReader([]byte(notification.DiffHTML))
	doc := tgbotapi.FileReader{
		Name:   "diff.html",
		Reader: diffBytes,
	}

	msg := tgbotapi.NewDocument(t.chatID, doc)
	msg.Caption = caption
	msg.ParseMode = "HTML"

	_, err := t.bot.Send(msg)
	if err != nil {
		return fmt.Errorf("failed to send telegram message: %w", err)
	}

	return nil
}
