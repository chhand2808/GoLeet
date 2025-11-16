package data

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type HistoryEntry struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Date  string `json:"date"` // YYYY-MM-DD
}

func (s *Store) HistoryPathInit() string {
	if s.HistoryPath == "" {
		s.HistoryPath = filepath.Join("data", "history.json")
	}
	return s.HistoryPath
}

// LoadHistory returns a slice of history entries (most recent last).
// If file doesn't exist or is empty/invalid, it will create/write an empty array and return empty slice.
func (s *Store) LoadHistory() ([]HistoryEntry, error) {
	hPath := s.HistoryPathInit()

	// create empty file if not exist
	if _, err := os.Stat(hPath); os.IsNotExist(err) {
		if err := os.WriteFile(hPath, []byte("[]"), 0644); err != nil {
			return nil, err
		}
	}

	raw, err := os.ReadFile(hPath)
	if err != nil {
		return nil, err
	}
	if len(raw) == 0 {
		return []HistoryEntry{}, nil
	}

	var hist []HistoryEntry
	if err := json.Unmarshal(raw, &hist); err == nil {
		return hist, nil
	}

	// if file contains object or corrupted, overwrite with empty array
	_ = s.SaveHistory([]HistoryEntry{})
	return []HistoryEntry{}, fmt.Errorf("history.json was invalid; reset to empty history")
}

// SaveHistory writes the history slice to disk (overwrites).
func (s *Store) SaveHistory(hist []HistoryEntry) error {
	hPath := s.HistoryPathInit()

	out, err := json.MarshalIndent(hist, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(hPath, out, 0644)
}

// AppendHistory appends a new entry to history as a queue (keeps maxLen entries).
// new entry will be appended to the end (most recent at end). If length exceeds maxLen,
// oldest entries are removed.
func (s *Store) AppendHistory(entry HistoryEntry, maxLen int) error {
	hist, err := s.LoadHistory()
	if err != nil && hist == nil {
		// even if error, ensure we have an empty slice to continue
		hist = []HistoryEntry{}
	}

	// Append entry
	hist = append(hist, entry)

	// Enforce max length
	if len(hist) > maxLen {
		// remove oldest entries from the front
		start := len(hist) - maxLen
		hist = hist[start:]
	}

	return s.SaveHistory(hist)
}

// Helper: create a HistoryEntry for a problem; date is set to today if empty
func NewHistoryEntry(id, title, date string) HistoryEntry {
	if date == "" {
		date = time.Now().Format("2006-01-02")
	}
	return HistoryEntry{
		ID:    id,
		Title: title,
		Date:  date,
	}
}
