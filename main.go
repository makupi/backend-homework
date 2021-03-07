package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/makupi/backend-homework/middlewares"
	"github.com/makupi/backend-homework/models"
	"github.com/makupi/backend-homework/storage"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

// getEnv returns content env if set, fallback if not
func getEnv(env, fallback string) string {
	if value, ok := os.LookupEnv(env); ok {
		return value
	}
	return fallback
}

// App contains the apps storage and JWTSecret
type App struct {
	Storage   storage.Storage
	JWTSecret []byte
}

// Initialize initializes the app with storage and loads the Secret
func (a *App) Initialize() {
	a.Storage = storage.NewSqliteStorage()
	a.JWTSecret = []byte(getEnv("JWT_SECRET", "development-secret"))
}

func addJSONPayload(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(payload)
	if err != nil {
		log.Print(err)
	}
}

func parseVarFromRequest(w http.ResponseWriter, r *http.Request, key string) (int, error) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars[key])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return id, err
	}
	return id, nil
}

// ListQuestions is the handler for GET /questions
func (a *App) ListQuestions(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int)
	lastID, _ := strconv.Atoi(r.URL.Query().Get("last_id"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	questions := a.Storage.List(userID, lastID, limit)
	addJSONPayload(w, http.StatusOK, questions)
}

// GetQuestion is the handler for GET /questions/{id}
func (a *App) GetQuestion(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int)
	id, err := parseVarFromRequest(w, r, "id")
	if err != nil {
		return
	}
	question, err := a.Storage.Get(id, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	addJSONPayload(w, http.StatusOK, question)
}

// UpdateQuestion is the handler for PUT /questions/{id}
func (a *App) UpdateQuestion(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int)
	id, err := parseVarFromRequest(w, r, "id")
	if err != nil {
		return
	}
	var question models.Question
	err = json.NewDecoder(r.Body).Decode(&question)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	question, err = a.Storage.Update(id, userID, question)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	addJSONPayload(w, http.StatusOK, question)
}

// NewQuestion is the handler for POST /questions
func (a *App) NewQuestion(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int)
	var question models.Question
	err := json.NewDecoder(r.Body).Decode(&question)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	question, err = a.Storage.Add(userID, question)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	addJSONPayload(w, http.StatusOK, question)
}

// DeleteQuestion is the handler for DELETE /questions/{id}
func (a *App) DeleteQuestion(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int)
	id, err := parseVarFromRequest(w, r, "id")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = a.Storage.Delete(id, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// AddOption is the handler for POST /questions/{id}/options
func (a *App) AddOption(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int)
	questionID, err := parseVarFromRequest(w, r, "id")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var option models.Option
	err = json.NewDecoder(r.Body).Decode(&option)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	question, err := a.Storage.AddOption(option, questionID, userID)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	}
	addJSONPayload(w, http.StatusOK, question)
}

// UpdateOption is the handler for PUT /questions/{id}/options/{id}
func (a *App) UpdateOption(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int)
	questionID, err := parseVarFromRequest(w, r, "id")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	optionID, err := parseVarFromRequest(w, r, "optionID")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var option models.Option
	err = json.NewDecoder(r.Body).Decode(&option)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	question, err := a.Storage.UpdateOption(option, optionID, questionID, userID)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	}
	addJSONPayload(w, http.StatusOK, question)
}

// DeleteOption is the handler for DELETE /questions/{id}/options/{id}
func (a *App) DeleteOption(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int)
	questionID, err := parseVarFromRequest(w, r, "id")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	optionID, err := parseVarFromRequest(w, r, "optionID")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var option models.Option
	err = json.NewDecoder(r.Body).Decode(&option)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	question, err := a.Storage.DeleteOption(optionID, questionID, userID)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	}
	addJSONPayload(w, http.StatusOK, question)
}

// CreateUser is the handler for POST /users
func (a *App) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	response, err := a.Storage.CreateUser(user.Username, user.Password)
	if err != nil {
		http.Error(w, "username already in use", http.StatusBadRequest)
		return
	}
	addJSONPayload(w, http.StatusOK, response)
}

// CreateToken is the handler for POST /users/token
func (a *App) CreateToken(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	token, err := a.Storage.CreateToken(user.Username, user.Password, a.JWTSecret)
	if err != nil {
		http.Error(w, "user does not exist or wrong password", http.StatusBadRequest)
		return
	}
	addJSONPayload(w, http.StatusOK, token)
}

func main() {
	app := App{}
	app.Initialize()
	jwtMiddleware := middlewares.JWTMiddleware{Secret: app.JWTSecret, Storage: app.Storage}
	router := mux.NewRouter()
	router.Use(middlewares.LoggingMiddleware)

	questions := router.PathPrefix("/questions").Subrouter()
	questions.Use(jwtMiddleware.Middleware)
	questions.HandleFunc("", app.ListQuestions).Methods("GET")
	questions.HandleFunc("", app.NewQuestion).Methods("POST")
	questions.HandleFunc("/{id}", app.GetQuestion).Methods("GET")
	questions.HandleFunc("/{id}", app.UpdateQuestion).Methods("PUT")
	questions.HandleFunc("/{id}", app.DeleteQuestion).Methods("DELETE")
	questions.HandleFunc("/{id}/options", app.AddOption).Methods("POST")
	questions.HandleFunc("/{id}/options/{optionID}", app.UpdateOption).Methods("PUT")
	questions.HandleFunc("/{id}/options/{optionID}", app.DeleteOption).Methods("DELETE")

	users := router.PathPrefix("/users").Subrouter()
	users.HandleFunc("", app.CreateUser).Methods("POST")
	users.HandleFunc("/token", app.CreateToken).Methods("POST")

	server := &http.Server{
		Addr:         getEnv("HOST", "127.0.0.1") + ":" + getEnv("PORT", "3000"),
		Handler:      router,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}
	log.Print("Running on " + server.Addr)
	log.Fatal(server.ListenAndServe())
}
