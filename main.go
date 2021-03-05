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

// getPort returns the content of the environment variable PORT if set, else 3000
func getPort() string {
	if port, ok := os.LookupEnv("PORT"); ok {
		return port
	}
	return "3000"
}

type App struct {
	Storage storage.Storage
}

func (a *App) Initialize() {
	a.Storage = storage.NewSqliteStorage()
}

func addJSONPayload(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(payload)
	if err != nil {
		log.Fatal(err)
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
	questions := a.Storage.List()
	addJSONPayload(w, http.StatusOK, questions)
}

func (a *App) GetQuestion(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDFromRequest(w, r)
	if err != nil {
		return
	}
	question, err := a.Storage.Get(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	addJSONPayload(w, http.StatusOK, question)
}

func (a *App) UpdateQuestion(w http.ResponseWriter, r *http.Request) {
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
	question, err = a.Storage.Update(id, question)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	addJSONPayload(w, http.StatusOK, question)
}

func (a *App) NewQuestion(w http.ResponseWriter, r *http.Request) {
	var question models.Question
	err := json.NewDecoder(r.Body).Decode(&question)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	question, err = a.Storage.Add(question)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	addJSONPayload(w, http.StatusOK, question)
}

func (a *App) DeleteQuestion(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDFromRequest(w, r)
	if err != nil {
		return
	}
	err = a.Storage.Delete(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func main() {
	app := App{}
	app.Initialize()
	router := mux.NewRouter()
	router.Use(middlewares.LoggingMiddleware)
	router.HandleFunc("/questions", app.ListQuestions).Methods("GET")
	router.HandleFunc("/questions/{id}", app.GetQuestion).Methods("GET")
	router.HandleFunc("/questions/{id}", app.UpdateQuestion).Methods("PUT")
	router.HandleFunc("/questions", app.NewQuestion).Methods("POST")
	router.HandleFunc("/questions/{id}", app.DeleteQuestion).Methods("DELETE")

	server := &http.Server{
		Addr:         "127.0.0.1:" + getPort(),
		Handler:      router,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}
	log.Fatal(server.ListenAndServe())
}
