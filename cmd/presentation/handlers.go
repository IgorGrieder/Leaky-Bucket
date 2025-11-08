package presentation

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/IgorGrieder/Leaky-Bucket/internal/application"
	"github.com/IgorGrieder/Leaky-Bucket/internal/domain"
)

type MutationHandler func(w http.ResponseWriter, r *http.Request, user *domain.User)

type MutationResponse struct {
	Data []domain.Mutation `json:"data"`
}

// @Summary Process Leaky Bucket Mutation
// @Description Attempts to process a mutation request, checking rate limits via the Leaky Bucket algorithm.
// @Tags Gateway
// @Accept json
// @Produce json
// @Param request body domain.Mutation true "Mutation Request Data"
// @Success 200 {object} MutationResponse "Request processed successfully"
// @Failure 429 {string} string "Rate limit exceeded (Leaky Bucket)"
// @Failure 500 {object} string "Internal server error"
// @Router /mutation [post]
// @security BearerAuth
func NewMutationHandler(service application.ProcessorService) MutationHandler {

	return func(w http.ResponseWriter, r *http.Request, user *domain.User) {
		var request domain.Mutation
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, "invalid request object", http.StatusBadRequest)
			return
		}

		if len(request.PIX_KEY) == 0 {
			http.Error(w, "invalid request object", http.StatusBadRequest)
			return
		}

		pix_keys, err := service.ProcessMutation(request, r.Context(), user)
		if err != nil {

			if errors.Is(err, &application.NoTokensError{}) {
				http.Error(w, "too many requests", http.StatusTooManyRequests)
				return
			}

			http.Error(w, "an error occured", http.StatusInternalServerError)
			return
		}

		if len(pix_keys) == 0 {
			http.Error(w, "no matches found for the request key", http.StatusNotFound)
			return
		}

		response := &MutationResponse{Data: pix_keys}

		returnJson, err := json.Marshal(response)
		if err != nil {
			http.Error(w, "failed to create JSON response", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(returnJson)
	}
}

// @Summary Generating JWT Token
// @Description Returns a valid JWT Token
// @Tags Gateway
// @Accept json
// @Produce json
// @Param request body domain.User true "User"
// @Success 204
// @Failure 400 {string} string "invalid request object"
// @Failure 500 {string} string "failed generating jwt token"
// @Router /generateJWT [post]
func Authenticate(authService application.AuthService) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		var request *domain.User
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, "invalid request object", http.StatusBadRequest)
			return
		}

		if len(request.Id) == 0 {
			http.Error(w, "invalid request object", http.StatusBadRequest)
			return
		}

		token, err := authService.GenerateToken(request.Id)
		if err != nil {
			log.Printf("error generating token %v", err)
			http.Error(w, "failed generating jwt token", http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "JWT_Token",
			Value:    token,
			HttpOnly: true,
			SameSite: http.SameSiteDefaultMode,
			Expires:  time.Now().Add(1 * time.Hour),
		})
		w.WriteHeader(http.StatusNoContent)
	}
}
