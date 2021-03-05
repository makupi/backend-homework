package main

type Question struct {
	Body    string   `json:"body"`
	Options []Option `json:"options"`
}

type Option struct {
	Body    string `json:"body"`
	Correct bool   `json:"correct"`
}
