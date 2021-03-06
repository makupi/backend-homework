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

type App struct {
	Storage   storage.Storage
	JWTSecret []byte
}

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

func parseIDFromRequest(w http.ResponseWriter, r *http.Request) (int, error) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return id, err
	}
	return id, nil
}

func (a *App) ListQuestions(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int)
	lastID, _ := strconv.Atoi(r.URL.Query().Get("last_id"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	questions := a.Storage.List(userID, lastID, limit)
	addJSONPayload(w, http.StatusOK, questions)
}

func (a *App) GetQuestion(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int)
	id, err := parseIDFromRequest(w, r)
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

func (a *App) UpdateQuestion(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int)
	id, err := parseIDFromRequest(w, r)
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

func (a *App) DeleteQuestion(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int)
	id, err := parseIDFromRequest(w, r)
	if err != nil {
		return
	}
	err = a.Storage.Delete(id, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func main() {
	app := App{}
	app.Initialize()
	jwtMiddleware := middlewares.JWTMiddleware{Secret: app.JWTSecret, Storage: app.Storage}
	router := mux.NewRouter()
	router.Use(jwtMiddleware.Middleware)
	router.Use(middlewares.LoggingMiddleware)
	router.HandleFunc("/questions", app.ListQuestions).Methods("GET")
	router.HandleFunc("/questions/{id}", app.GetQuestion).Methods("GET")
	router.HandleFunc("/questions/{id}", app.UpdateQuestion).Methods("PUT")
	router.HandleFunc("/questions", app.NewQuestion).Methods("POST")
	router.HandleFunc("/questions/{id}", app.DeleteQuestion).Methods("DELETE")

	server := &http.Server{
		Addr:         "127.0.0.1:" + getEnv("PORT", "3000"),
		Handler:      router,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}
	log.Print("Running on " + server.Addr)
	log.Fatal(server.ListenAndServe())
}
