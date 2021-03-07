package middlewares

import (
	"log"
	"net/http"
)

// Simple logging middleware that logs "<Method> <URI>" before executing the request
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
	})
}
