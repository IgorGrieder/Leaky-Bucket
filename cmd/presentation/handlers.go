package presentation

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/IgorGrieder/Leaky-Bucket/internal/application"
	"github.com/IgorGrieder/Leaky-Bucket/internal/domain"
)

type MutationHandler func(w http.ResponseWriter, r *http.Request, user *domain.User)

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

		pix_keys, err := service.ProcessMutation(request, r.Context())
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

		returnJson, err := json.Marshal(pix_keys)
		if err != nil {
			http.Error(w, "failed to create JSON response", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(returnJson)
	}
}

func Authenticate(authService application.AuthService) http.HandlerFunc {
	type response struct {
		Token string `json:"token"`
	}

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
		response := &response{Token: token}

		if err != nil {
			log.Printf("error generating token %v", err)
			http.Error(w, "failed generating jwt token", http.StatusBadRequest)
			return
		}

		returnJson, err := json.Marshal(response)
		if err != nil {
			log.Printf("error generating token %v", err)
			http.Error(w, "failed generating jwt token", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(returnJson)
	}

}
