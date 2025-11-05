package presentation

import (
	"net/http"
	"strings"

	"github.com/IgorGrieder/Leaky-Bucket/internal/application"
	"github.com/IgorGrieder/Leaky-Bucket/internal/config"
	"github.com/IgorGrieder/Leaky-Bucket/internal/ctx"
)

func AuthMiddleware(next http.HandlerFunc, cfg *config.Config) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "missing authorization header", http.StatusUnauthorized)
			return
		}

		token := trimBearerPrefix(tokenString)
		if token == tokenString {
			http.Error(w, "invalid authorization header format", http.StatusUnauthorized)
			return
		}

		tokenParsed, err := application.Authenticate(token, cfg.HASH)
		if err != nil || !tokenParsed.Valid {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		claims, ok := tokenParsed.Claims.(*application.JWT)
		if !ok {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		newCtx := ctx.SetUserIdCtx(r.Context(), claims.UserID)
		next(w, r.WithContext(newCtx))
	})
}

func trimBearerPrefix(token string) string {
	return strings.TrimPrefix(token, "Bearer ")
}
