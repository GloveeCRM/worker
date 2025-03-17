package main

import (
	"context"
	"fmt"
	"glovee-worker/config"
	"glovee-worker/database"
	"glovee-worker/service/email"
	"glovee-worker/types"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// Interface for handling all task types
type TaskHandler interface {
	HandleTask(task *types.Task) (map[string]any, error)
}

type EmailTaskHandler struct {
	emailService *email.Service
}

func NewEmailTaskHandler(config *types.Config) *EmailTaskHandler {
	return &EmailTaskHandler{
		emailService: email.NewService(config),
	}
}

func (h *EmailTaskHandler) HandleTask(task *types.Task) (map[string]any, error) {
	emailData := task.Data.(types.Email)
	err := h.emailService.SendEmail(&emailData)

	metadata := map[string]any{
		"email_id": emailData.EmailID,
	}

	return metadata, err
}

// TODO: Add other communication channels here
// E.g. SMS, Push Notifications, etc.

func main() {
	config := config.NewConfig()

	db, err := database.NewDB(context.Background(), config)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize task handlers
	taskHandlers := map[types.TaskType]TaskHandler{
		types.TaskTypeEmail: NewEmailTaskHandler(config),
		// TODO: Add other communication channels here
		// E.g. types.TaskTypeSMS: NewSMSTaskHandler(config),
	}

	// Create a context that can be cancelled
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Set up signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start the worker loop in a goroutine
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				// Process tasks for each task type
				for taskType, handler := range taskHandlers {
					if err := processTask(ctx, db, taskType, handler); err != nil {
						log.Printf("Error processing task type %s: %v", taskType, err)
					}
				}
				time.Sleep(5 * time.Second)
			}
		}
	}()

	// Wait for shutdown signal
	<-sigChan
	fmt.Println("\nShutting down gracefully...")
	cancel()
	time.Sleep(1 * time.Second) // Give ongoing tasks a chance to complete
}

func processTask(ctx context.Context, db *database.DB, taskType types.TaskType, handler TaskHandler) error {
	task, err := db.DequeueTask(ctx, taskType)
	if err != nil {
		return fmt.Errorf("failed to dequeue task: %w", err)
	}

	if task == nil {
		log.Printf("No task available for task type %s", taskType)
		return nil // No task available
	}

	metadata, err := handler.HandleTask(task)

	var errorMessage string
	if err != nil {
		errorMessage = err.Error()
		log.Printf("Failed to process task type %s: %v", taskType, err)
	}

	result, err := db.ProcessTaskResult(ctx, task.TaskID, err == nil, errorMessage, metadata)
	if err != nil {
		return fmt.Errorf("failed to process task result: %w", err)
	}

	if result.Success {
		log.Printf("Task type %s completed successfully", taskType)
	} else {
		log.Printf("Task type %s failed: %s", taskType, result.Error)
	}

	return nil
}
