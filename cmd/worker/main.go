package main

import (
	"context"
	"fmt"
	"glovee-worker/config"
	"glovee-worker/database"
	"glovee-worker/service/email"
	"glovee-worker/types"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	config := config.NewConfig()

	db, err := database.NewDB(context.Background(), config)
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

	emailService := email.NewService(config)
	emailData := task.Data.(types.Email)
	err = emailService.SendEmail(&emailData)

	metadata := map[string]any{
		"email_id": emailData.EmailID,
	}

	var errorMessage string
	if err != nil {
		errorMessage = err.Error()
		log.Printf("Failed to send email: %v", err)
	}

	result, err := db.ProcessTaskResult(context.Background(), task.TaskID, err == nil, errorMessage, metadata)
	if err != nil {
		log.Fatalf("Failed to process task result: %v", err)
	}

	if result.Success {
		fmt.Println("Email sent and task completed successfully")
	} else {
		fmt.Printf("Task failed: %s\n", result.Error)
	}
}
