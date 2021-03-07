package models

// Question is the JSON representation for questions over the REST API
type Question struct {
	ID      int      `json:"id"`
	Body    string   `json:"body"`
	Options []Option `json:"options"`
}

// Option is the JSON representation for options over the REST API
type Option struct {
	ID         int    `json:"id"`
	Body       string `json:"body"`
	Correct    bool   `json:"correct"`
	QuestionID int    `json:"-"`
}
