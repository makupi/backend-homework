package middlewares

import (
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/makupi/backend-homework/models"
	"github.com/makupi/backend-homework/storage"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// JWTMiddleware holds the secret and storage needed for the JWT Middleware
type JWTMiddleware struct {
	Secret  []byte
	Storage storage.Storage
}

// Middleware that checks for a JWT token and verifies if the userID inside exists
// "userID" will be set in r.Context
func (j *JWTMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		bearer := strings.TrimPrefix(auth, "Bearer ")
		token, err := jwt.Parse(bearer, func(token *jwt.Token) (interface{}, error) {
			return j.Secret, nil
		})
		if err == nil && token.Valid {
			userID := (*token).Claims.(jwt.MapClaims)["userID"]
			if userID != nil {
				userID, err := strconv.Atoi(fmt.Sprintf("%v", userID))
				if err == nil && j.Storage.UserIDExists(userID) {
					ctx := context.WithValue(r.Context(), models.ContextUserID, userID)
					next.ServeHTTP(w, r.WithContext(ctx))
					return
				}
			}
		}
		log.Print("Unauthorized access to " + r.Method + " " + r.RequestURI)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	})
}
