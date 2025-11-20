package gemini

import "strings"

// ExtractJSON isolates the JSON array inside any Gemini text output.
// It finds the first '[' and last ']' and returns the substring.
func ExtractJSON(raw string) string {
	start := strings.Index(raw, "[")
	end := strings.LastIndex(raw, "]")

	if start == -1 || end == -1 || start >= end {
		return "" // no valid JSON array found
	}

	return raw[start : end+1]
}
