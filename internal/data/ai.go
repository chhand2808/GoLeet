package data

type AISuggestion struct {
	Title  string   `json:"title"`
	Number int      `json:"number"`
	Topics []string `json:"topics"`
}
