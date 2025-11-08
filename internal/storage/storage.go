package storage

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const (
	// StateFile is the JSON file tracking endpoint hashes
	StateFile = "jsmon.json"
	// DownloadsDir is where fetched content is stored
	DownloadsDir = "downloads"
)

// EndpointState maps endpoint URLs to their hash history
type EndpointState map[string][]string

// Storage handles reading and writing endpoint state and downloaded files
type Storage struct {
	stateFile    string
	downloadsDir string
}

// New creates a new Storage instance
func New(stateFile, downloadsDir string) *Storage {
	return &Storage{
		stateFile:    stateFile,
		downloadsDir: downloadsDir,
	}
}

// NewDefault creates a Storage instance with default paths
func NewDefault() *Storage {
	return New(StateFile, DownloadsDir)
}

// Initialize ensures storage directories and files exist
func (s *Storage) Initialize() error {
	// Create downloads directory if it doesn't exist
	if err := os.MkdirAll(s.downloadsDir, 0755); err != nil {
		return fmt.Errorf("failed to create downloads directory: %w", err)
	}

	// Create state file if it doesn't exist
	if _, err := os.Stat(s.stateFile); os.IsNotExist(err) {
		emptyState := make(EndpointState)
		if err := s.SaveState(emptyState); err != nil {
			return fmt.Errorf("failed to create state file: %w", err)
		}
	}

	return nil
}

// LoadState reads the endpoint state from jsmon.json
func (s *Storage) LoadState() (EndpointState, error) {
	data, err := os.ReadFile(s.stateFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read state file: %w", err)
	}

	var state EndpointState
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, fmt.Errorf("failed to parse state file: %w", err)
	}

	return state, nil
}

// SaveState writes the endpoint state to jsmon.json
func (s *Storage) SaveState(state EndpointState) error {
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal state: %w", err)
	}

	if err := os.WriteFile(s.stateFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write state file: %w", err)
	}

	return nil
}

// GetPreviousHash returns the most recent hash for an endpoint, or empty string if none exists
func (s *Storage) GetPreviousHash(state EndpointState, endpoint string) string {
	hashes, exists := state[endpoint]
	if !exists || len(hashes) == 0 {
		return ""
	}
	return hashes[len(hashes)-1]
}

// AddHash appends a new hash to an endpoint's history
func (s *Storage) AddHash(state EndpointState, endpoint, hash string) {
	state[endpoint] = append(state[endpoint], hash)
}

// SaveEndpointContent saves the fetched content to a file named by its hash
func (s *Storage) SaveEndpointContent(hash, content string) error {
	filePath := filepath.Join(s.downloadsDir, hash)
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to save endpoint content: %w", err)
	}
	return nil
}

// LoadEndpointContent loads content from a hash-named file
func (s *Storage) LoadEndpointContent(hash string) (string, error) {
	filePath := filepath.Join(s.downloadsDir, hash)
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to load endpoint content: %w", err)
	}
	return string(data), nil
}

// GetFileSize returns the size of a downloaded file in bytes
func (s *Storage) GetFileSize(hash string) (int64, error) {
	filePath := filepath.Join(s.downloadsDir, hash)
	info, err := os.Stat(filePath)
	if err != nil {
		return 0, fmt.Errorf("failed to stat file: %w", err)
	}
	return info.Size(), nil
}

// HashContent calculates MD5 hash of content (first 10 chars like Python version)
func HashContent(content string) string {
	hash := md5.Sum([]byte(content))
	return fmt.Sprintf("%x", hash)[:10]
}

// GetDownloadPath returns the full path for a hash file
func (s *Storage) GetDownloadPath(hash string) string {
	return filepath.Join(s.downloadsDir, hash)
}
