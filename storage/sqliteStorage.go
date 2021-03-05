package storage

import (
	"database/sql"
	"github.com/togglhire/backend-homework"
)

type SqliteStorage struct {
	DB *sql.DB
}

func (s *SqliteStorage) List() []main.Question {

}

func (s *SqliteStorage) Add(question main.Question) (main.Question, error) {

}

func (s *SqliteStorage) Get(id int) (main.Question, error) {

}

func (s *SqliteStorage) Update(id int, question main.Question) (main.Question, error) {

}