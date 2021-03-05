package main

import (
	"github.com/gorilla/mux"
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

func ListQuestions(w http.ResponseWriter, r *http.Request) {

}

func GetQuestion(w http.ResponseWriter, r *http.Request) {

}

func UpdateQuestion(w http.ResponseWriter, r *http.Request) {

}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/questions", ListQuestions).Methods("GET")
	router.HandleFunc("/questions/{id}", GetQuestion).Methods("GET")
	router.HandleFunc("/questions/{id}", UpdateQuestion).Methods("PUT")

	server := &http.Server {
		Addr: "127.0.0.1:" + getPort(),
		Handler: router,
		ReadTimeout: 1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}
	log.Fatal(server.ListenAndServe())
}
