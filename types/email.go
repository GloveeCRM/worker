package types

import "time"

type EmailStatusType string

const (
	EmailStatusTypePending EmailStatusType = "pending"
	EmailStatusTypeSent    EmailStatusType = "sent"
	EmailStatusTypeFailed  EmailStatusType = "failed"
)

type EmailStatus struct {
	StatusID  int64           `json:"status_id"`
	EmailID   int64           `json:"email_id"`
	Status    EmailStatusType `json:"status"`
	Attempts  int             `json:"attempts"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

type Email struct {
	EmailID   int64       `json:"email_id"`
	FromEmail string      `json:"from_email"`
	ToEmail   string      `json:"to_email"`
	Subject   string      `json:"subject"`
	HTML      string      `json:"html"`
	CreatedAt time.Time   `json:"created_at"`
	Status    EmailStatus `json:"status"`
}
