package notifier

import (
	"fmt"

	"github.com/slack-go/slack"
)

// SlackNotifier sends notifications via Slack
type SlackNotifier struct {
	client    *slack.Client
	channelID string
}

// NewSlackNotifier creates a new Slack notifier
func NewSlackNotifier(token, channelID string) *SlackNotifier {
	return &SlackNotifier{
		client:    slack.New(token),
		channelID: channelID,
	}
}

// Name returns the notifier name
func (s *SlackNotifier) Name() string {
	return "Slack"
}

// Notify sends a change notification via Slack
func (s *SlackNotifier) Notify(notification *ChangeNotification) error {
	message := fmt.Sprintf(
		"[JSMon] %s has been updated! Download the diff HTML file below to check changes.",
		notification.Endpoint,
	)

	params := slack.FileUploadParameters{
		Channels:       []string{s.channelID},
		Content:        notification.DiffHTML,
		Filename:       "diff.html",
		Filetype:       "html",
		Title:          "Diff Changes",
		InitialComment: message,
	}

	_, err := s.client.UploadFile(params)
	if err != nil {
		return fmt.Errorf("failed to upload file to slack: %w", err)
	}

	return nil
}
