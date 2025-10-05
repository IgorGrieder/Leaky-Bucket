package presentation

import (
	"context"
	"net/http"
	"strings"

	"github.com/IgorGrieder/Leaky-Bucket/internal/application"
	"github.com/IgorGrieder/Leaky-Bucket/internal/config"
)

func AuthMiddleware(next http.HandlerFunc, cfg *config.Config) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "Missing authorization header", http.StatusUnauthorized)
			return
		}

		token := trimBearerPrefix(tokenString)
		if token == tokenString {
			http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
			return
		}

		tokenParsed, err := application.Authenticate(token, cfg.HASH)
		if err != nil || !tokenParsed.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "JWT", tokenParsed)
		next(w, r.WithContext(ctx))
	})
}

func trimBearerPrefix(token string) string {
	return strings.TrimPrefix(token, "Bearer ")
}
