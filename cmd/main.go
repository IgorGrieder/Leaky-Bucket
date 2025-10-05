package main

import (
	"fmt"
	"os"

	"github.com/IgorGrieder/Leaky-Bucket/cmd/presentation"
	"github.com/IgorGrieder/Leaky-Bucket/internal/application"
	"github.com/IgorGrieder/Leaky-Bucket/internal/config"
	"github.com/IgorGrieder/Leaky-Bucket/internal/database"
	"github.com/IgorGrieder/Leaky-Bucket/internal/repository"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	cfg := config.NewConfig()

	if err != nil {
		fmt.Println("ENV loading failed, canceling the program")
		os.Exit(1)
	}

	fmt.Println("Starting the program")

	connections := database.StartConns(cfg)

	LimitingRepository := repository.NewLimitingRepository(connections.Redis)
	MutationRepository := repository.NewMutationRepository(connections.PG)

	gatewayService := &application.ProcessorService{
		MutationRepository: MutationRepository,
		LimitingRepository: LimitingRepository,
	}

	fmt.Println("Root layer stablished, starting the http server")

	presentation.StartHttpServer(cfg, gatewayService)
}
