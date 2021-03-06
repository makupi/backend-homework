package storage

import (
	"database/sql"
	"github.com/makupi/backend-homework/models"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

type SqliteStorage struct {
	DB *sql.DB
}

func NewSqliteStorage() *SqliteStorage {
	db, err := sql.Open("sqlite3", "./db.sqlite3?_foreign_keys=on")
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
		log.Print(err)
	}
	defer rows.Close()

	for rows.Next() {
		var option models.Option
		if err := rows.Scan(&option.ID, &option.QuestionID, &option.Body, &option.Correct); err != nil {
			log.Print(err)
		}
		options = append(options, option)
	}
	return
}

func (s *SqliteStorage) List(lastID, limit int) (questions []models.Question) {
	var rows *sql.Rows
	var err error
	if (lastID != 0) && (limit != 0) {
		rows, err = s.DB.Query(
			`SELECT * FROM questions WHERE id < (?) ORDER BY ID DESC LIMIT (?)`,
			lastID,
			limit,
		)
		if err != nil {
			log.Print(err)
		}
	} else {
		rows, err = s.DB.Query(`SELECT * FROM questions`)
		if err != nil {
			log.Print(err)
		}
	}

	defer rows.Close()

	for rows.Next() {
		var question models.Question
		if err := rows.Scan(&question.ID, &question.Body); err != nil {
			log.Print(err)
		}
		question.Options = s.getOptions(question.ID)
		questions = append(questions, question)
	}
	return
}

func (s *SqliteStorage) addOptions(options []models.Option, questionID int) error {
	for _, option := range options {
		_, err := s.DB.Exec(
			`INSERT INTO options (QUESTION_ID, OPTION, CORRECT) values (?,?,?)`,
			questionID,
			option.Body,
			option.Correct,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *SqliteStorage) Add(question models.Question) (models.Question, error) {
	var q models.Question
	result, err := s.DB.Exec(`INSERT INTO questions (question) values (?)`, question.Body)
	if err != nil {
		return q, nil
	}
	id, err := result.LastInsertId()
	if err != nil {
		return q, nil
	}
	err = s.addOptions(question.Options, int(id))
	return s.Get(int(id))
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

func (s *SqliteStorage) updateQuestion(id int, question models.Question) error {
	_, err := s.DB.Exec(`UPDATE questions SET question = (?) WHERE ID == (?)`, question.Body, id)
	return err
}

func (s *SqliteStorage) updateOption(option models.Option) error {
	_, err := s.DB.Exec(
		`UPDATE options SET OPTION = (?), CORRECT = (?) WHERE ID == (?)`,
		option.Body,
		option.Correct,
		option.ID,
	)
	return err
}

func (s *SqliteStorage) Update(id int, question models.Question) (models.Question, error) {
	currentQ, err := s.Get(id)
	if err != nil {
		return models.Question{}, err
	}
	if currentQ.Body != question.Body {
		err = s.updateQuestion(id, question)
		if err != nil {
			return models.Question{}, err
		}
	}
	for _, option := range question.Options {
		for _, currentOption := range currentQ.Options {
			if option.ID == currentOption.ID {
				if (option.Body != currentOption.Body) || (option.Correct != currentOption.Correct) {
					err = s.updateOption(option)
					if err != nil {
						return models.Question{}, err
					}
				}
			}
		}
	}
	return s.Get(id)
}

func (s *SqliteStorage) Delete(id int) error {
	_, err := s.DB.Exec(`DELETE FROM questions WHERE ID == (?)`, id)
	return err
}

func (s *SqliteStorage) UserIDExists(userID int) bool {
	return true
}
