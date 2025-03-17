package email

import (
	"fmt"
	"glovee-worker/types"

	"github.com/resend/resend-go/v2"
)

type Service struct {
	client *resend.Client
}

func NewService(config *types.Config) *Service {
	return &Service{
		client: resend.NewClient(config.Email.ResendAPIKey),
	}
}

func (s *Service) SendEmail(email *types.Email) error {
	const operation = "service.email.SendEmail"

	emailRequest := resend.SendEmailRequest{
		From:    email.FromEmail,
		To:      []string{email.ToEmail},
		Subject: email.Subject,
		Html:    email.HTML,
	}

	_, err := s.client.Emails.Send(&emailRequest)
	if err != nil {
		return fmt.Errorf("%s: %w", operation, err)
	}

	return nil
}
