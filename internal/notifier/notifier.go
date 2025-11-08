package notifier

import (
	"fmt"
)

// ChangeNotification contains information about a detected change
type ChangeNotification struct {
	Endpoint    string
	OldHash     string
	NewHash     string
	OldSize     int64
	NewSize     int64
	DiffHTML    string
	DiffContent []byte
}

// Notifier is the interface all notification backends must implement
type Notifier interface {
	Notify(notification *ChangeNotification) error
	Name() string
}

// MultiNotifier sends notifications to multiple backends
type MultiNotifier struct {
	notifiers []Notifier
}

// NewMultiNotifier creates a notifier that sends to multiple backends
func NewMultiNotifier(notifiers ...Notifier) *MultiNotifier {
	return &MultiNotifier{
		notifiers: notifiers,
	}
}

// Notify sends the notification to all configured backends
func (m *MultiNotifier) Notify(notification *ChangeNotification) error {
	var errors []error

	for _, notifier := range m.notifiers {
		if err := notifier.Notify(notification); err != nil {
			errors = append(errors, fmt.Errorf("%s: %w", notifier.Name(), err))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("notification errors: %v", errors)
	}

	return nil
}

// FormatMessage creates a standardized notification message
func FormatMessage(notification *ChangeNotification) string {
	return fmt.Sprintf(
		"%s has been updated from %s (%d bytes) to %s (%d bytes)",
		notification.Endpoint,
		notification.OldHash,
		notification.OldSize,
		notification.NewHash,
		notification.NewSize,
	)
}
