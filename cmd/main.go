package main

import (
	"log"

	"github.com/IgorGrieder/Leaky-Bucket/cmd/presentation"
	"github.com/IgorGrieder/Leaky-Bucket/internal/application"
	"github.com/IgorGrieder/Leaky-Bucket/internal/config"
	"github.com/IgorGrieder/Leaky-Bucket/internal/database"
	"github.com/IgorGrieder/Leaky-Bucket/internal/repository"
)

func main() {
	log.Println("Starting the program")

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

	log.Println("root layer stablished, starting the http server")

	// HTTP Server
	presentation.StartHttpServer(cfg, gatewayService, authService)
}
