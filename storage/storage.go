package storage

import "github.com/togglhire/backend-homework"

type Storage interface {
	List() []main.Question
	Add(question main.Question) (main.Question, error)
	Get(id int) (main.Question, error)
	Update(id int, question main.Question) (main.Question, error)
}