package main

import (
	"context"
	"fmt"
	"glovee-worker/internal/config"
	"glovee-worker/internal/repository/postgres"
	"log"

	_ "github.com/lib/pq"
)

func main() {
	config := config.NewConfig()

	db, err := postgres.NewDB(context.Background(), config)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	hi, err := db.SayHello()
	if err != nil {
		log.Fatalf("Failed to say hello: %v", err)
	}
	fmt.Println(hi)
}
