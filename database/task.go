package database

import (
	"context"
	"encoding/json"
	"fmt"
	"glovee-worker/types"
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
	dataBytes, err := json.Marshal(task.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal task data: %w", err)
	}

	switch task.TaskType {
	case types.TaskTypeEmail:
		var emailData struct {
			Email types.Email `json:"email"`
		}
		if err := json.Unmarshal(dataBytes, &emailData); err != nil {
			return nil, fmt.Errorf("failed to unmarshal email data: %w", err)
		}
		task.Data = emailData.Email
		// TODO: Add other task types here
		// E.g. types.TaskTypeSMS:
		// var smsData struct {
		// 	SMS types.SMS `json:"sms"`
		// }
		// if err := json.Unmarshal(dataBytes, &smsData); err != nil {
		// 	return nil, fmt.Errorf("failed to unmarshal sms data: %w", err)
		// }
	}

	return &task, nil
}

func (db *DB) ProcessTaskResult(ctx context.Context, taskID int64, success bool, errorMessage string, metadata any) (*types.TaskResult, error) {
	var result []byte
	metadataJSON, err := json.Marshal(metadata)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal metadata: %w", err)
	}

	err = db.pool.QueryRow(ctx, `SELECT queues.process_task_result($1, $2, $3, $4)`,
		taskID, success, errorMessage, metadataJSON).Scan(&result)
	if err != nil {
		return nil, fmt.Errorf("failed to process task result: %w", err)
	}

	var taskResult types.TaskResult
	if err := json.Unmarshal(result, &taskResult); err != nil {
		return nil, fmt.Errorf("failed to unmarshal task result: %w", err)
	}

	return &taskResult, nil
}
