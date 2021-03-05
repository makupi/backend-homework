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

func (s *SqliteStorage) getOptions(questionID int) (options []models.Option) {
	rows, err := s.DB.Query(`SELECT * FROM options WHERE QUESTION_ID == (?)`, questionID)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var option models.Option
		if err := rows.Scan(&option.ID, &option.QuestionID, &option.Body, &option.Correct); err != nil {
			log.Fatal(err)
		}
		options = append(options, option)
	}
	return
}

func (s *SqliteStorage) List() (questions []models.Question) {
	rows, err := s.DB.Query(`SELECT * FROM questions`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var question models.Question
		if err := rows.Scan(&question.ID, &question.Body); err != nil {
			log.Fatal(err)
		}
		question.Options = s.getOptions(question.ID)
		questions = append(questions, question)
	}
	return
}

func (s *SqliteStorage) Add(question models.Question) (models.Question, error) {
	var q models.Question
	return q, errors.New("not implemented")
}

func (s *SqliteStorage) Get(id int) (models.Question, error) {
	row := s.DB.QueryRow(`SELECT * FROM questions WHERE ID == (?)`, id)
	var question models.Question
	err := row.Scan(&question.ID, &question.Body)
	if err != nil {
		return question, err
	}
	question.Options = s.getOptions(question.ID)
	return question, nil
}

func (s *SqliteStorage) Update(id int, question models.Question) (models.Question, error) {
	var q models.Question
	return q, errors.New("not implemented")
}
