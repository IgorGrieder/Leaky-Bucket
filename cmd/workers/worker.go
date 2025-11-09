package workers

import (
	"log/slog"
	"time"

	"github.com/IgorGrieder/Leaky-Bucket/internal/application"
)

func TokenRefillWorker(service application.ProcessorService) {
	timer := time.NewTicker(1 * time.Minute)
	defer timer.Stop()

	for range timer.C {

		err := service.FetchAndRefilTokens()
		if err != nil {
			slog.Error("token refill job failed", slog.Any("error", err))
		}

		slog.Info("token refill job succeeded")
	}
}
