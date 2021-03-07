package models

// JSON representation of questions
type Question struct {
	ID      int      `json:"id"`
	Body    string   `json:"body"`
	Options []Option `json:"options"`
}

// JSON representation of options
type Option struct {
	ID         int    `json:"id"`
	Body       string `json:"body"`
	Correct    bool   `json:"correct"`
	QuestionID int    `json:"-"`
}
