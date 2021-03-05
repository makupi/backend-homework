package models

type Question struct {
	ID      int      `json:"id"`
	Body    string   `json:"body"`
	Options []Option `json:"options"`
}

type Option struct {
	ID         int    `json:"id"`
	Body       string `json:"body"`
	Correct    bool   `json:"correct"`
	QuestionID int    `json:"-"`
}
