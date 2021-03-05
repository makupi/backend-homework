package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/makupi/backend-homework/storage"
	"log"
	"net/http"
	"os"
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

func (a *App) ListQuestions(w http.ResponseWriter, r *http.Request) {
	questions := a.Storage.List()
	addJSONPayload(w, http.StatusOK, questions)
	fmt.Println(questions)
}

func (a *App) GetQuestion(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	fmt.Println("GET /questions/" + id)
}

func (a *App) UpdateQuestion(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	fmt.Println("PUT /questions/" + id)
}

func main() {
	app := App{}
	app.Initialize()
	router := mux.NewRouter()
	router.HandleFunc("/questions", app.ListQuestions).Methods("GET")
	router.HandleFunc("/questions/{id}", app.GetQuestion).Methods("GET")
	router.HandleFunc("/questions/{id}", app.UpdateQuestion).Methods("PUT")

	server := &http.Server{
		Addr:         "127.0.0.1:" + getPort(),
		Handler:      router,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}
	log.Fatal(server.ListenAndServe())
}
