package main

import (
	"context"
	"encoding/json"
	"fmt"
	"glovee-worker/internal/config"
	"glovee-worker/internal/repository/postgres"
	"glovee-worker/internal/types"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	config := config.NewConfig()

	db, err := postgres.NewDB(context.Background(), config)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	task, err := db.DequeueTask(context.Background(), types.TaskTypeEmail)
	if err != nil {
		log.Fatalf("Failed to dequeue task: %v", err)
	}

	if task == nil {
		fmt.Println("No task available")
		return
	}

	jsonBytes, err := json.MarshalIndent(task, "", "    ")
	if err != nil {
		log.Fatalf("Failed to marshal task to JSON: %v", err)
	}
	fmt.Println(string(jsonBytes))
}
