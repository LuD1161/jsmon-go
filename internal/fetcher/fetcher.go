package fetcher

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"
)

// TargetsDir is the default directory for target files
const TargetsDir = "targets"

// Fetcher handles fetching endpoints
type Fetcher struct {
	client     *http.Client
	maxRetries int
}

// New creates a new Fetcher with custom settings
func New(timeout time.Duration, maxRetries int) *Fetcher {
	return &Fetcher{
		client: &http.Client{
			Timeout: timeout,
		},
		maxRetries: maxRetries,
	}
}

// NewDefault creates a Fetcher with default settings
func NewDefault() *Fetcher {
	return New(30*time.Second, 3)
}

// LoadEndpoints reads all endpoint URLs from files in the targets directory
func LoadEndpoints(targetsDir string) ([]string, error) {
	var endpoints []string

	entries, err := os.ReadDir(targetsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read targets directory: %w", err)
	}

	for _, entry := range entries {
		// Skip hidden files and directories
		if entry.IsDir() || strings.HasPrefix(entry.Name(), ".") {
			continue
		}

		filePath := filepath.Join(targetsDir, entry.Name())
		data, err := os.ReadFile(filePath)
		if err != nil {
			return nil, fmt.Errorf("failed to read target file %s: %w", entry.Name(), err)
		}

		// Split by lines and trim whitespace
		lines := strings.Split(string(data), "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" || strings.HasPrefix(line, "#") {
				continue // Skip empty lines and comments
			}
			if IsValidURL(line) {
				endpoints = append(endpoints, line)
			}
		}
	}

	return endpoints, nil
}

// IsValidURL validates if a string is a valid HTTP/HTTPS/FTP URL
func IsValidURL(url string) bool {
	// Regex from Python version
	pattern := `^(?:http|ftp)s?://(?:(?:[A-Z0-9](?:[A-Z0-9-]{0,61}[A-Z0-9])?\.)+(?:[A-Z]{2,6}\.?|[A-Z0-9-]{2,}\.?)|localhost|\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})(?::\d+)?(?:/?|[/?]\S+)$`
	matched, _ := regexp.MatchString("(?i)"+pattern, url)
	return matched
}

// FetchResult holds the result of fetching a single endpoint
type FetchResult struct {
	Endpoint string
	Content  string
	Error    error
}

// FetchAll fetches all endpoints concurrently
func (f *Fetcher) FetchAll(endpoints []string) []FetchResult {
	var wg sync.WaitGroup
	results := make([]FetchResult, len(endpoints))

	for i, endpoint := range endpoints {
		wg.Add(1)
		go func(idx int, ep string) {
			defer wg.Done()
			content, err := f.Fetch(ep)
			results[idx] = FetchResult{
				Endpoint: ep,
				Content:  content,
				Error:    err,
			}
		}(i, endpoint)
	}

	wg.Wait()
	return results
}

// Fetch retrieves content from a single endpoint with retry logic
func (f *Fetcher) Fetch(endpoint string) (string, error) {
	var lastErr error

	for attempt := 0; attempt <= f.maxRetries; attempt++ {
		if attempt > 0 {
			// Exponential backoff
			time.Sleep(time.Duration(attempt) * time.Second)
		}

		resp, err := f.client.Get(endpoint)
		if err != nil {
			lastErr = err
			continue
		}

		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			lastErr = fmt.Errorf("HTTP %d: %s", resp.StatusCode, resp.Status)
			continue
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			lastErr = err
			continue
		}

		return string(body), nil
	}

	return "", fmt.Errorf("failed after %d attempts: %w", f.maxRetries+1, lastErr)
}
