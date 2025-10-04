package main

import (
	"fmt"
	"os"

	"github.com/IgorGrieder/Leaky-Bucket/internal/database"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("ENV loading failed, canceling the program")
		os.Exit(1)
	}

	fmt.Println("Starting the program")

	database.StartConns()

	fmt.Println("Connections stablished")

}
