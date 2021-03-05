package storage

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/togglhire/backend-homework/models"
	"log"
)

type SqliteStorage struct {
	DB *sql.DB
}

func NewSqliteStorage() *SqliteStorage {
	db, err := sql.Open("sqlite3", "./db.sqlite3")
	if err != nil {
		log.Fatal(err)
	}
	storage := SqliteStorage{DB: db}
	return &storage
}

func (s *SqliteStorage) List() []models.Question {

}

func (s *SqliteStorage) Add(question models.Question) (models.Question, error) {

}

func (s *SqliteStorage) Get(id int) (models.Question, error) {

}

func (s *SqliteStorage) Update(id int, question models.Question) (models.Question, error) {

}
