package storage

import (
	"github.com/makupi/backend-homework/models"
)

type Storage interface {
	List(userID, lastID, limit int) []models.Question
	Add(userID int, question models.Question) (models.Question, error)
	Get(id, userID int) (models.Question, error)
	Update(id, userID int, question models.Question) (models.Question, error)
	Delete(id, userID int) error
	CreateUser(username, password string) error
	//CreateToken(username, password string) (string, error)
	UserIDExists(userID int) bool
}
