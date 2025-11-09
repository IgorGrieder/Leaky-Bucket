package main

import (
	"log/slog"
	"os"

	"github.com/IgorGrieder/Leaky-Bucket/cmd/presentation"
	"github.com/IgorGrieder/Leaky-Bucket/cmd/workers"
	"github.com/IgorGrieder/Leaky-Bucket/internal/application"
	"github.com/IgorGrieder/Leaky-Bucket/internal/config"
	"github.com/IgorGrieder/Leaky-Bucket/internal/database"
	"github.com/IgorGrieder/Leaky-Bucket/internal/repository"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	slog.Info("Starting the program")

	// ENVs
	cfg := config.NewConfig()

	// Database connections
	connections := database.StartConns(cfg)

	// Repositories
	LimitingRepository := repository.NewLimitingRepository(connections.Redis)
	MutationRepository := repository.NewMutationRepository(connections.PG)

	// Services
	gatewayService := application.ProcessorService{
		MutationRepository: MutationRepository,
		LimitingRepository: LimitingRepository,
	}

	authService := application.NewAuthService(cfg)

	slog.Info("root layer stablished, starting the refill worker")
	go workers.TokenRefillWorker(gatewayService)

	slog.Info("root layer stablished, starting the http server")
	// HTTP Server
	presentation.StartHttpServer(cfg, gatewayService, authService)
}
