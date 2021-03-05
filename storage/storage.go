package storage

import (
	"github.com/togglhire/backend-homework/models"
)

type Storage interface {
	List() []models.Question
	Add(question models.Question) (models.Question, error)
	Get(id int) (models.Question, error)
	Update(id int, question models.Question) (models.Question, error)
}
