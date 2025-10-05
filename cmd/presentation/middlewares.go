package presentation

import (
	"fmt"
	"net/http"
)

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// --- Authentication Logic Here ---
		// e.g., check header, validate token

		// If authorization fails:
		// http.Error(w, "Unauthorized", http.StatusUnauthorized)
		// return // Stop processing and don't call next

		// If authorization succeeds:
		// 2. Correct call: next(w, r)
		token := r.Header.Get("Authorization")
		if token == "" {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, "Missing authorization header")
			return
		}
		next(w, r)
	})
}
