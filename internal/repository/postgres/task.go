package postgres

import (
	"context"
	"encoding/json"
	"fmt"
	"glovee-worker/internal/types"
)

func (db *DB) DequeueTask(ctx context.Context, taskType types.TaskType) (*types.Task, error) {
	var result []byte
	err := db.pool.QueryRow(ctx, `SELECT queues.dequeue_task($1)`, taskType).Scan(&result)
	if err != nil {
		return nil, fmt.Errorf("failed to dequeue task: %w", err)
	}

	// If no task was dequeued
	if result == nil {
		return nil, nil
	}

	var task types.Task
	if err := json.Unmarshal(result, &task); err != nil {
		return nil, fmt.Errorf("failed to unmarshal task: %w", err)
	}

	// Handle the data field based on task type
	if task.TaskType == types.TaskTypeEmail {
		var emailData struct {
			Email types.Email `json:"email"`
		}
		dataBytes, err := json.Marshal(task.Data)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal task data: %w", err)
		}
		if err := json.Unmarshal(dataBytes, &emailData); err != nil {
			return nil, fmt.Errorf("failed to unmarshal email data: %w", err)
		}
		task.Data = emailData.Email
	}

	return &task, nil
}
