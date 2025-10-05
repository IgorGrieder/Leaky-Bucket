package main

import (
	"fmt"
	"os"

	"github.com/IgorGrieder/Leaky-Bucket/cmd/presentation"
	"github.com/IgorGrieder/Leaky-Bucket/internal/application"
	"github.com/IgorGrieder/Leaky-Bucket/internal/config"
	"github.com/IgorGrieder/Leaky-Bucket/internal/database"
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
	gatewayService := application.ProcessorService{}

	fmt.Println("Connections stablished")

	presentation.StartHttpServer(cfg, gatewayService)
}
