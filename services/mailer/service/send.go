package service

import (
	"net/smtp"

	"github.com/google/uuid"

	"github.com/pkg/errors"

	"github.com/artem-malko/auth-and-go/services/mailer"
)

func (s *mailerService) SendEmail(priority mailer.EmailPriority, email mailer.Email) error {
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	message := `To: "` + email.Name + `" <` + email.Email + `>
From: "Example" <info@example.com>
Subject: ` + email.Subject + "\n" + mime + "\n" + email.Message

	err := smtp.SendMail(s.host+":"+s.port, s.auth, s.from, []string{email.Email}, []byte(message))

	return errors.Wrap(err, "mailer service: SendEmail error")
}

func (s *mailerService) SendRegistrationConfirmationEmail(email string, confirmationToken uuid.UUID) error {
	templateData := struct {
		Email string
		URL   string
	}{
		Email: email,
		URL:   "https://dev.example.com/api/confirm?token=" + confirmationToken.String(),
	}
	body, err := s.parseTemplate("services/mailer/service/templates/confirm_registration.html", templateData)

	if err != nil {
		return errors.Wrap(err, "mailer service: SendRegistrationConfirmationEmail")
	}

	err = s.SendEmail(mailer.EmailPriorityHigh, mailer.Email{
		Email:   email,
		Name:    email,
		Subject: "Confirm registration",
		Message: body,
	})

	if err != nil {
		return errors.Wrap(err, "mailer service: SendRegistrationConfirmationEmail")
	}

	return nil
}

func (s *mailerService) SendUnexpectedRegistrationEmail(email, name, username string) error {
	nameForEmail := name

	if name == "" {
		nameForEmail = username
	}

	templateData := struct {
		Email string
		Name  string
	}{
		Email: email,
		Name:  nameForEmail,
	}
	body, err := s.parseTemplate("services/mailer/service/templates/unexpected_registration.html", templateData)
	if err != nil {
		return errors.Wrap(err, "mailer service: SendUnexpectedRegistrationEmail")
	}

	err = s.SendEmail(mailer.EmailPriorityHigh, mailer.Email{
		Email:   email,
		Name:    email,
		Subject: "Unexpected registration",
		Message: body,
	})

	if err != nil {
		return errors.Wrap(err, "mailer service: SendUnexpectedRegistrationEmail")
	}

	return nil
}

func (s *mailerService) SendRegistrationConfirmedEmail(email, name, username string) error {
	nameForEmail := name

	if name == "" {
		nameForEmail = username
	}

	templateData := struct {
		Name string
	}{
		Name: nameForEmail,
	}
	body, err := s.parseTemplate("services/mailer/service/templates/registration_confirmed.html", templateData)
	if err != nil {
		return errors.Wrap(err, "mailer service: SendRegistrationConfirmedEmail")
	}

	err = s.SendEmail(mailer.EmailPriorityHigh, mailer.Email{
		Email:   email,
		Name:    email,
		Subject: "Registration confirmed, welcome!",
		Message: body,
	})

	if err != nil {
		return errors.Wrap(err, "mailer service: SendRegistrationConfirmedEmail")
	}

	return nil
}
