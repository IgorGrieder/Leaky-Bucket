package main

import (
	"fmt"

	"github.com/IgorGrieder/Leaky-Bucket/internal/database"
)

func main() {
	fmt.Println("Starting the program")

	database.StartConns()

	fmt.Println("Connections stablished")

}
