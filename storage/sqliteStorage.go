package storage

import (
	"database/sql"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/makupi/backend-homework/models"
	_ "github.com/mattn/go-sqlite3" // driver for sqlite3
	"log"
)

// SqliteStorage object to access database
type SqliteStorage struct {
	DB *sql.DB
}

// NewSqliteStorage Creates a local db.sqlite3 database and automaticlly creates tables
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

// createTables creates tables necessary database tables
// users:
// | id: pkey, int | username: text, unique | password: text |
// questions:
// | id: pkey, int | body: text | userID: fkey(users.id), int |
// options:
// | id: pkey, int | body: text | correct: bool | question_id: fkey(questions.id), int |
func (s *SqliteStorage) createTables() error {
	tables := []string{
		`CREATE TABLE IF NOT EXISTS "users" (
			"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			"username" TEXT NOT NULL UNIQUE,
			"password" TEXT
		);`,
		`CREATE TABLE IF NOT EXISTS "questions" (
			"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			"question" TEXT,
			"user_id" INTEGER NOT NULL,
			CONSTRAINT fk_user_id
				FOREIGN KEY (user_id)
				REFERENCES users(id)
				ON DELETE CASCADE
		);`,
		`CREATE TABLE IF NOT EXISTS "options" (
			"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			"question_id" INTEGER NOT NULL,
			"option" TEXT,
			"correct" BOOLEAN,
			CONSTRAINT fk_question_id
				FOREIGN KEY (question_id)
				REFERENCES questions(id)
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
	rows, err := s.DB.Query(`SELECT * FROM options WHERE question_id == (?)`, questionID)
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

// List returns all questions that belong to the userID
func (s *SqliteStorage) List(userID, lastID, limit int) (questions []models.Question) {
	var rows *sql.Rows
	var err error
	if (lastID != 0) && (limit != 0) {
		rows, err = s.DB.Query(
			`SELECT * FROM questions WHERE user_id == (?) AND id < (?) ORDER BY id DESC LIMIT (?)`,
			userID,
			lastID,
			limit,
		)
		if err != nil {
			log.Print(err)
		}
	} else {
		rows, err = s.DB.Query(`SELECT * FROM questions WHERE user_id == (?)`, userID)
		if err != nil {
			log.Print(err)
		}
	}

	defer rows.Close()
	questions = []models.Question{}
	for rows.Next() {
		var question models.Question
		var _userID int
		if err := rows.Scan(&question.ID, &question.Body, &_userID); err != nil {
			log.Print(err)
		}
		question.Options = s.getOptions(question.ID)
		questions = append(questions, question)
	}
	return
}

// AddOption adds an Option to an existing question
// If the question doesn't belong to userID it will result in an error
func (s *SqliteStorage) AddOption(option models.Option, questionID, userID int) (models.Question, error) {
	var question models.Question
	if !s.HasQuestionAccess(userID, questionID) {
		return question, fmt.Errorf("unauthorized")
	}
	_, err := s.DB.Exec(
		`INSERT INTO options (question_id, option, correct) values (?,?,?)`,
		questionID,
		option.Body,
		option.Correct,
	)
	if err != nil {
		return question, err
	}
	return s.Get(questionID, userID)
}

func (s *SqliteStorage) addOptions(options []models.Option, questionID int) error {
	for _, option := range options {
		_, err := s.DB.Exec(
			`INSERT INTO options (question_id, option, correct) values (?,?,?)`,
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

// Add a new Question associated to the userID
func (s *SqliteStorage) Add(userID int, question models.Question) (models.Question, error) {
	var q models.Question
	result, err := s.DB.Exec(`INSERT INTO questions (question, user_id) values (?, ?)`, question.Body, userID)
	if err != nil {
		return q, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return q, err
	}
	err = s.addOptions(question.Options, int(id))
	if err != nil {
		return q, err
	}
	return s.Get(int(id), userID)
}

// Get a question by ID, will only return questions associated to the userID
func (s *SqliteStorage) Get(id, userID int) (models.Question, error) {
	row := s.DB.QueryRow(`SELECT * FROM questions WHERE id == (?) AND user_id == (?)`, id, userID)
	var question models.Question
	var _userID int
	err := row.Scan(&question.ID, &question.Body, &_userID)
	if err != nil {
		return question, err
	}
	question.Options = s.getOptions(question.ID)
	return question, nil
}

func (s *SqliteStorage) updateQuestion(id, userID int, question models.Question) error {
	_, err := s.DB.Exec(`UPDATE questions SET question = (?) WHERE id == (?) AND user_id == (?)`, question.Body, id, userID)
	return err
}

// UpdateOption updates an existing option
// If the question doesn't belong to userID it will result in an error
func (s *SqliteStorage) UpdateOption(option models.Option, optionID, questionID, userID int) (models.Question, error) {
	var question models.Question
	if !s.HasQuestionAccess(userID, questionID) {
		return question, fmt.Errorf("unauthorized")
	}
	_, err := s.DB.Exec(
		`UPDATE options SET option = (?), correct = (?) WHERE id == (?) AND question_id == (?)`,
		option.Body,
		option.Correct,
		optionID,
		questionID,
	)
	if err != nil {
		return question, err
	}
	return s.Get(questionID, userID)
}

// Update updates an existing question
// If the question doesn't belong to userID it will result in an error
func (s *SqliteStorage) Update(id, userID int, question models.Question) (models.Question, error) {
	currentQ, err := s.Get(id, userID)
	if err != nil {
		return models.Question{}, err
	}
	if currentQ.Body != question.Body {
		err = s.updateQuestion(id, userID, question)
		if err != nil {
			return models.Question{}, err
		}
	}
	for _, option := range question.Options {
		for _, currentOption := range currentQ.Options {
			if option.ID == currentOption.ID {
				if (option.Body != currentOption.Body) || (option.Correct != currentOption.Correct) {
					_, err := s.DB.Exec(
						`UPDATE options SET option = (?), correct = (?) WHERE id == (?) AND question_id == (?)`,
						option.Body,
						option.Correct,
						option.ID,
						id,
					)
					if err != nil {
						return models.Question{}, err
					}
				}
			}
		}
	}
	return s.Get(id, userID)
}

// DeleteOption deletes an existing option from a question
//If the question doesn't belong to userID or it doesn't exist it will result in an error
func (s *SqliteStorage) DeleteOption(optionID, questionID, userID int) (models.Question, error) {
	var question models.Question
	if !s.HasQuestionAccess(userID, questionID) {
		return question, fmt.Errorf("unauthorized")
	}
	_, err := s.DB.Exec(`DELETE FROM options WHERE id == (?) AND question_id == (?)`, optionID, questionID)
	if err != nil {
		return question, err
	}
	return s.Get(questionID, userID)
}

// Delete deletes an existing question
//If the question doesn't belong to userID or it doesn't exist it will result in an error
func (s *SqliteStorage) Delete(id, userID int) error {
	_, err := s.DB.Exec(`DELETE FROM questions WHERE id == (?) AND user_id == (?)`, id, userID)
	return err
}

// CreateUser creates a new user with username and password
// If the user already exists it will return an error
func (s *SqliteStorage) CreateUser(username, password string) (models.UserResponse, error) {
	result, err := s.DB.Exec(`INSERT INTO users (username, password) values (?, ?)`, username, password)
	if err != nil {
		return models.UserResponse{}, err
	}
	id, err := result.LastInsertId()
	return models.UserResponse{ID: int(id), Username: username}, err
}

// CreateToken creates a new JWT token for the user
// If username and password are incorrect it will result in an error
func (s *SqliteStorage) CreateToken(username, password string, secret []byte) (models.JWTTokenResponse, error) {
	var jwtToken models.JWTTokenResponse
	var user models.User
	row := s.DB.QueryRow(`SELECT * FROM users WHERE username == (?) AND password == (?)`, username, password)
	err := row.Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		return jwtToken, err
	}
	token := jwt.New(jwt.GetSigningMethod("HS256"))
	(*token).Claims.(jwt.MapClaims)["userID"] = user.ID
	tokenString, err := token.SignedString(secret)
	if err != nil {
		return jwtToken, err
	}
	jwtToken.Token = tokenString

	return jwtToken, nil
}

// UserIDExists checks if a given userID exists
func (s *SqliteStorage) UserIDExists(userID int) bool {
	row := s.DB.QueryRow(`SELECT * FROM users WHERE id == (?)`, userID)
	var user models.User
	err := row.Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		return false
	}
	return true
}

// HasQuestionAccess verifies that a userID has access to a questionID
// Returns true if the user has access and false if not
func (s *SqliteStorage) HasQuestionAccess(userID, questionID int) bool {
	row := s.DB.QueryRow(`SELECT questions.id FROM questions WHERE ID == (?) AND user_id == (?)`, questionID, userID)
	var question models.Question
	err := row.Scan(&question.ID)
	if err != nil {
		return false
	}
	return true
}
