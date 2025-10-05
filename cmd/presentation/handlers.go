package presentation

import (
	"encoding/json"
	"net/http"

	"github.com/IgorGrieder/Leaky-Bucket/internal/application"
	"github.com/IgorGrieder/Leaky-Bucket/internal/domain"
)

func MutationHandler(w http.ResponseWriter, r *http.Request) {
	var request domain.Mutation
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request object", http.StatusBadRequest)
		return
	}

	application.ProcessMutation()

	//I need to call a service here

}
