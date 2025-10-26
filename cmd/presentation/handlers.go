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

		if len(request.PIX_KEY) == 0 {
			http.Error(w, "Invalid request object", http.StatusBadRequest)
			return
		}

		pix_keys, err := service.ProcessMutation(request, r.Context())
		if err != nil {
			http.Error(w, "An error occured", http.StatusInternalServerError)
			return
		}

		if len(pix_keys) == 0 {
			http.Error(w, "No matches found for the request key", http.StatusNotFound)
			return
		}

		returnJson, err := json.Marshal(pix_keys)
		if err != nil {
			http.Error(w, "Failed to create JSON response", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(returnJson)
	}
}

func Authenticate(w http.ResponseWriter, r *http.Request) {

}
