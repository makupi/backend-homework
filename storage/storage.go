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
	CreateUser(username, password string) (models.UserResponse, error)
	CreateToken(username, password string, secret []byte) (models.JWTTokenResponse, error)
	UserIDExists(userID int) bool
	HasQuestionAccess(userID, questionID int) bool
	AddOption(option models.Option, questionID, userID int) (models.Question, error)
	UpdateOption(option models.Option, optionID, questionID, userID int) (models.Question, error)
	DeleteOption(optionID, questionID, userID int) (models.Question, error)
}
