package presentation

import (
	"encoding/json"
	"net/http"

	"github.com/IgorGrieder/Leaky-Bucket/internal/application"
	"github.com/IgorGrieder/Leaky-Bucket/internal/domain"
)

func NewMutationHandler(service application.ProcessorService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request domain.Mutation
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, "Invalid request object", http.StatusBadRequest)
			return
		}

		// If we get an error we must return a error code
		pix_key, err := service.ProcessMutation(request, r.Context())
		if err != nil {
			http.Error(w, "An error occured", http.StatusInternalServerError)
			return
		}

		returnJson, err := json.Marshal(pix_key)
		if err != nil {
			http.Error(w, "Failed to create JSON response", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(returnJson)
	}
}
