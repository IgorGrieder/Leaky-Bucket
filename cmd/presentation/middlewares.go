package presentation

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/IgorGrieder/Leaky-Bucket/internal/config"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(next func(w http.ResponseWriter, r *http.Request, ctx context.Context), cfg *config.Config) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, "Missing authorization header")
			return
		}

		tokenStringTreated := strings.TrimPrefix(tokenString, "Bearer ")
		if tokenStringTreated == tokenString {
			http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
			return
		}

		token, err := jwt.Parse(tokenStringTreated, func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			return cfg.HASH, nil
		})

		if err != nil {
			http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
			return
		}

		if !token.Valid {
			http.Error(w, "Expired token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "JWT", token)
		next(w, r, ctx)
	})
}
