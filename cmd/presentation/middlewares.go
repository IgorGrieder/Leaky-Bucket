package presentation

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/IgorGrieder/Leaky-Bucket/internal/application"
	"github.com/IgorGrieder/Leaky-Bucket/internal/config"
	"github.com/IgorGrieder/Leaky-Bucket/internal/domain"
)

func AuthMiddleware(handler MutationHandler, authService application.AuthService, cfg *config.Config) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		cookie, err := r.Cookie("JWT_Token")
		if err != nil {
			if err == http.ErrNoCookie {
				http.Error(w, "JWT cookie not found", http.StatusUnauthorized)
				return
			}
			http.Error(w, "Failed to retrieve cookie", http.StatusInternalServerError)
			return
		}

		tokenString := cookie.Value
		if tokenString == "" {
			http.Error(w, "missing authorization header", http.StatusUnauthorized)
			return
		}

		token := trimBearerPrefix(tokenString)
		if token == tokenString {
			http.Error(w, "invalid authorization header format", http.StatusUnauthorized)
			return
		}

		tokenParsed, err := authService.Authenticate(token, cfg.HASH)
		if err != nil || !tokenParsed.Valid {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		claims, ok := tokenParsed.Claims.(*application.JWT)
		if !ok {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		user := &domain.User{Id: claims.UserID}

		slog.Info("user authenticated", slog.Group("user_credentials", slog.String("user_id", user.Id)))

		handler(w, r, user)
	})
}

func trimBearerPrefix(token string) string {
	return strings.TrimPrefix(token, "Bearer ")
}
