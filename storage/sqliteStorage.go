package storage

import (
	"database/sql"
	"errors"
	"github.com/makupi/backend-homework/models"
	_ "github.com/mattn/go-sqlite3"
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
	err = storage.createTables()
	if err != nil {
		log.Fatal(err)
	}
	return &storage
}

// creates tables questions and options
// questions:
// | id: pkey, int | body: text |
// options:
// | id: pkey, int | body: text | correct: bool | question_id: fkey(questions.id), int |
func (s *SqliteStorage) createTables() error {
	tables := []string{
		`CREATE TABLE IF NOT EXISTS "questions" (
		"ID" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		"question" TEXT
		);`,
		`CREATE TABLE IF NOT EXISTS "options" (
			"ID" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			"QUESTION_ID" INTEGER NOT NULL,
			"OPTION" TEXT,
			"CORRECT" BOOLEAN,
			CONSTRAINT fk_questions
				FOREIGN KEY (QUESTION_ID)
				REFERENCES questions(ID)
				ON DELETE CASCADE
		);`,
	}

	for _, table := range tables {
		_, err := s.DB.Exec(table)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *SqliteStorage) List() []models.Question {
	return nil
}

func (s *SqliteStorage) Add(question models.Question) (models.Question, error) {
	var q models.Question
	return q, errors.New("not implemented")
}

func (s *SqliteStorage) Get(id int) (models.Question, error) {
	var q models.Question
	return q, errors.New("not implemented")
}

func (s *SqliteStorage) Update(id int, question models.Question) (models.Question, error) {
	var q models.Question
	return q, errors.New("not implemented")
}
