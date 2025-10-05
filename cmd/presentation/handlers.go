package presentation

import (
	"encoding/json"
	"net/http"

	"github.com/IgorGrieder/Leaky-Bucket/internal/application"
	"github.com/IgorGrieder/Leaky-Bucket/internal/domain"
)

func NewMutationHandler(service *application.ProcessorService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request domain.Mutation
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, "Invalid request object", http.StatusBadRequest)
			return
		}

		if err := service.ProcessMutation(request, r.Context()); err != nil {
			http.Error(w, "An error occured", http.StatusInternalServerError)
			return
		}
	}
}
