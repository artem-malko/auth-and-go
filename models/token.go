package models

import (
	"time"

	"github.com/google/uuid"
)

type TokenType string

func (t TokenType) String() string {
	return string(t)
}

const (
	TokenTypeRegistrationConfirmation TokenType = "registration_confirmation"
	TokenTypeEmailConfirmation        TokenType = "email_confirmation"
	TokenTypeAutoLogin                TokenType = "auto_login"
)

type TokenStatus string

func (t TokenStatus) String() string {
	return string(t)
}

const (
	TokenStatusActive TokenStatus = "active"
	TokenStatusUsed   TokenStatus = "used"
)

type Token struct {
	ID          uuid.UUID
	TokenType   TokenType
	TokenStatus TokenStatus
	AccountID   uuid.UUID
	IdentityID  uuid.UUID
	ClientID    ClientID
	ExpiresDate time.Time
}

func CheckTokenType(tokenType string) bool {
	switch tokenType {
	case TokenTypeRegistrationConfirmation.String():
		fallthrough
	case TokenTypeEmailConfirmation.String():
		fallthrough
	case TokenTypeAutoLogin.String():
		return true
	}

	return false
}
