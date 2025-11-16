package model

type TopicTag struct {
	Name string `json:"name"`
}

type Problem struct {
	ID         string     `json:"frontendQuestionId"`
	Title      string     `json:"title"`
	Slug       string     `json:"titleSlug"`
	Difficulty string     `json:"difficulty"`
	TopicTags  []TopicTag `json:"topicTags"`
}
