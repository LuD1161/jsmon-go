package notifier

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
)

// DiscordNotifier sends notifications via Discord webhooks
type DiscordNotifier struct {
	webhookURL string
	client     *http.Client
}

// NewDiscordNotifier creates a new Discord notifier
func NewDiscordNotifier(webhookURL string) *DiscordNotifier {
	return &DiscordNotifier{
		webhookURL: webhookURL,
		client:     &http.Client{},
	}
}

// Name returns the notifier name
func (d *DiscordNotifier) Name() string {
	return "Discord"
}

// Notify sends a change notification via Discord webhook
func (d *DiscordNotifier) Notify(notification *ChangeNotification) error {
	message := fmt.Sprintf(
		"%s has been updated from ***%s*** (***%d*** bytes) to ***%s*** (***%d*** bytes)",
		notification.Endpoint,
		notification.OldHash,
		notification.OldSize,
		notification.NewHash,
		notification.NewSize,
	)

	// Create multipart form data
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add content field
	if err := writer.WriteField("content", message); err != nil {
		return fmt.Errorf("failed to write content field: %w", err)
	}

	if err := writer.WriteField("username", "JSMon"); err != nil {
		return fmt.Errorf("failed to write username field: %w", err)
	}

	// Add file
	part, err := writer.CreateFormFile("file", "diff.html")
	if err != nil {
		return fmt.Errorf("failed to create form file: %w", err)
	}

	if _, err := io.WriteString(part, notification.DiffHTML); err != nil {
		return fmt.Errorf("failed to write file content: %w", err)
	}

	if err := writer.Close(); err != nil {
		return fmt.Errorf("failed to close multipart writer: %w", err)
	}

	// Send request
	req, err := http.NewRequest("POST", d.webhookURL, body)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := d.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send discord webhook: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("discord webhook returned status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	return nil
}
