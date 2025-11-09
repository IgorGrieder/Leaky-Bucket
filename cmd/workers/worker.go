package workers

import (
	"log"
	"time"

	"github.com/IgorGrieder/Leaky-Bucket/internal/application"
)

func TokenRefillWorker(service application.ProcessorService) {
	timer := time.NewTicker(1 * time.Minute)
	defer timer.Stop()

	for range timer.C {

		err := service.FetchAndRefilTokens()
		if err != nil {
			log.Printf("token refill job failed: %v", err)
		}

		log.Printf("token refill job succeeded")
	}
}
