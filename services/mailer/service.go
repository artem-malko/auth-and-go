package mailer

import "github.com/google/uuid"

type EmailPriority string

const (
	EmailPriorityHigh   EmailPriority = "high"
	EmailPriorityMedium EmailPriority = "medium"
	EmailPriorityLow    EmailPriority = "low"
)

type Email struct {
	Email   string
	Name    string
	Subject string
	Message string
}

type Service interface {
	SendEmail(priority EmailPriority, email Email) error
	SendRegistrationConfirmationEmail(email string, confirmationToken uuid.UUID) error
	SendUnexpectedRegistrationEmail(email, name, username string) error
	SendRegistrationConfirmedEmail(email, name, username string) error
}
