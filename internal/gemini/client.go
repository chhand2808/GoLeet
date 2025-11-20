package gemini

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/chhand2808/goleet/internal/data"
)

// Use Gemini 2.0 Flash-Lite model
const model = "models/gemini-2.0-flash-lite"

func GetSuggestions(prompt string) ([]data.AISuggestion, error) {
	apiKey := loadAPIKey()
	if apiKey == "" {
		return nil, errors.New("API key missing. Run `goleet init` again")
	}

	reqBody := map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"role": "user",
				"parts": []map[string]string{
					{"text": prompt},
				},
			},
		},
	}

	jsonData, _ := json.Marshal(reqBody)

	// REST endpoint for Gemini 2.0 Flash-Lite
	url := fmt.Sprintf(
		"https://generativelanguage.googleapis.com/v1beta/%s:generateContent?key=%s",
		model, apiKey,
	)

	resp, err := http.Post(url, "application/json", bytes.NewReader(jsonData))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("gemini API error (%d): %s", resp.StatusCode, string(body))
	}

	var gResp GeminiResponse
	if err := json.Unmarshal(body, &gResp); err != nil {
		return nil, fmt.Errorf("invalid Gemini response JSON: %w\nRaw: %s",
			err, string(body))
	}

	if len(gResp.Candidates) == 0 ||
		len(gResp.Candidates[0].Content.Parts) == 0 {
		return nil, errors.New("empty AI response")
	}

	text := gResp.Candidates[0].Content.Parts[0].Text

	// extract clean JSON array
	clean := ExtractJSON(text)
	if clean == "" {
		return nil, fmt.Errorf("AI output does not contain JSON array.\nRaw: %s", text)
	}

	var parsed []data.AISuggestion
	if err := json.Unmarshal([]byte(clean), &parsed); err != nil {
		return nil, fmt.Errorf("AI output JSON parse error: %v\nCleaned: %s", err, clean)
	}

	return parsed, nil
}

func loadAPIKey() string {
	b, err := os.ReadFile("data/config.json")
	if err != nil {
		return ""
	}

	var cfg map[string]string
	json.Unmarshal(b, &cfg)
	return cfg["api_key"]
}
