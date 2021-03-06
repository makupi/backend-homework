package middlewares

import (
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/makupi/backend-homework/storage"
	"log"
	"net/http"
	"strings"
)

type JWTMiddleware struct {
	Secret  []byte
	Storage storage.Storage
}

func (j *JWTMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		bearer := strings.TrimPrefix(auth, "Bearer ")
		token, err := jwt.Parse(bearer, func(token *jwt.Token) (interface{}, error) {
			return j.Secret, nil
		})
		if err == nil && token.Valid {
			userID := (*token).Claims.(jwt.MapClaims)["userID"]
			if userID != nil && j.Storage.UserIDExists(fmt.Sprintf("%v", userID)) {
				ctx := context.WithValue(r.Context(), "userID", userID)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}
		}
		log.Print("Unauthorized access to " + r.Method + " " + r.RequestURI)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	})
}
