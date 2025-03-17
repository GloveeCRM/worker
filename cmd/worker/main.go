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
	if err != nil {
		log.Fatalf("Failed to send email: %v", err)
	}

	fmt.Println("Email sent successfully")
}
