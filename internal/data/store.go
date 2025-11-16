package data

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type Problem struct {
	ID         string `json:"frontendQuestionId"`
	Title      string `json:"title"`
	Difficulty string `json:"difficulty"`
	TitleSlug  string `json:"titleSlug"`
	TopicTags  []struct {
		Name string `json:"name"`
	} `json:"topicTags"`
}

type SolvedProblem struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Date  string `json:"date"`
}

type Store struct {
	ProblemsPath string
	SolvedPath   string
	HistoryPath  string
}

func NewStore() *Store {
	return &Store{
		ProblemsPath: filepath.Join("data", "problems.json"),
		SolvedPath:   filepath.Join("data", "solved.json"),
		HistoryPath:  filepath.Join("data", "history.json"),
	}
}

// Load all problems
func (s *Store) LoadProblems() ([]Problem, error) {
	file, err := os.ReadFile(s.ProblemsPath)
	if err != nil {
		return nil, err
	}

	var problems []Problem
	err = json.Unmarshal(file, &problems)
	if err != nil {
		return nil, err
	}

	return problems, nil
}

// Load solved problems
func (s *Store) LoadSolved() ([]SolvedProblem, error) {
	// Create file if not exists
	if _, err := os.Stat(s.SolvedPath); os.IsNotExist(err) {
		err = os.WriteFile(s.SolvedPath, []byte("[]"), 0644)
		if err != nil {
			return nil, err
		}
	}

	file, err := os.ReadFile(s.SolvedPath)
	if err != nil {
		return nil, err
	}

	// If empty, return empty array
	if len(file) == 0 {
		return []SolvedProblem{}, nil
	}

	// Try unmarshalling into array
	var solved []SolvedProblem
	err = json.Unmarshal(file, &solved)
	if err == nil {
		return solved, nil
	}

	// If file was mistakenly an object {}, fix it
	var single map[string]interface{}
	err2 := json.Unmarshal(file, &single)
	if err2 == nil {
		// File contains `{}`, convert to empty list
		return []SolvedProblem{}, s.SaveSolved([]SolvedProblem{})
	}

	// JSON corrupted
	return nil, fmt.Errorf("solved.json is invalid; delete or fix the file: %v", err)
}

// Save solved problems
func (s *Store) SaveSolved(solved []SolvedProblem) error {
	data, err := json.MarshalIndent(solved, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.SolvedPath, data, 0644)
}

// Mark a problem as solved (adds or updates date)
func (s *Store) MarkSolved(problemID, title string) error {
	solved, err := s.LoadSolved()
	if err != nil {
		return err
	}

	date := time.Now().Format("2006-01-02")
	found := false

	for i, p := range solved {
		if p.ID == problemID {
			solved[i].Date = date
			found = true
			break
		}
	}

	if !found {
		solved = append(solved, SolvedProblem{
			ID:    problemID,
			Title: title,
			Date:  date,
		})
	}

	return s.SaveSolved(solved)
}
