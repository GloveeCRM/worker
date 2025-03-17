package types

import "time"

type TaskType string

const (
	TaskTypeEmail TaskType = "email"
)

type Task struct {
	TaskID       int64     `json:"task_id"`
	TaskType     TaskType  `json:"task_type"`
	ResourceID   int64     `json:"resource_id"`
	Priority     int       `json:"priority"`
	Retries      int       `json:"retries"`
	MaxRetries   int       `json:"max_retries"`
	CreatedAt    time.Time `json:"created_at"`
	ScheduledAt  time.Time `json:"scheduled_at"`
	DequeuedAt   time.Time `json:"dequeued_at"`
	CompletedAt  time.Time `json:"completed_at"`
	FailedAt     time.Time `json:"failed_at"`
	ErrorMessage string    `json:"error_message"`
	Metadata     any       `json:"metadata"`
	Data         any       `json:"data"`
}

type TaskResult struct {
	Success    bool   `json:"success"`
	Message    string `json:"message"`
	TaskType   string `json:"task_type"`
	ResourceID int64  `json:"resource_id"`
	Retries    int    `json:"retries,omitempty"`
	MaxRetries int    `json:"max_retries,omitempty"`
	Error      string `json:"error,omitempty"`
}
