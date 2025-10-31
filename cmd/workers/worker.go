package workers

import (
	"context"
	"log"
	"time"

	"github.com/IgorGrieder/Leaky-Bucket/internal/application"
)

func TokenRefillWorker(service application.ProcessorService) {
	timer := time.NewTicker(1 * time.Minute)
	defer timer.Stop()

	for range timer.C {
		ctx := context.Background()

		err := service.RefillTokens(ctx)
		if err != nil {
			log.Printf("token refill failed: %v", err)
		}

		log.Printf("token refill succeeded: %v", err)
	}
}
