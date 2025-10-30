package workers

import (
	"time"

	"github.com/IgorGrieder/Leaky-Bucket/internal/application"
)

func TokenRefillWorker(service application.ProcessorService) {
	tick := time.NewTicker(1 * time.Minute)
}
