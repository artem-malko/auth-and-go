package models

import (
	"errors"
	"html"
	"net/mail"
	"strings"
	"time"

	"github.com/google/uuid"
)

var addressParser = &mail.AddressParser{}

var (
	ErrIdentityEmptyEmail     = errors.New("identity email can't be empty")
	ErrIdentityIncorrectEmail = errors.New("identity email is incorrect")
	ErrIdentityIncorrectType  = errors.New("identity type is incorrect")
)

type IdentityType string

func (i IdentityType) String() string {
	return string(i)
}

type IdentityStatus string

func (i IdentityStatus) String() string {
	return string(i)
}

const (
	IdentityTypeEmail    IdentityType = "email"
	IdentityTypeFacebook IdentityType = "facebook"
	IdentityTypeGoogle   IdentityType = "google"
)

const (
	IdentityStatusUnconfirmed IdentityStatus = "unconfirmed"
	IdentityStatusConfirmed   IdentityStatus = "confirmed"
)

type SocialNetworkType string

const (
	SocialNetworkTypeGoogle   SocialNetworkType = "google"
	SocialNetworkTypeFacebook SocialNetworkType = "facebook"
)

type Identity struct {
	ID               uuid.UUID      `json:"-"`
	AccountID        uuid.UUID      `json:"-"`
	IdentityType     IdentityType   `json:"-"`
	IdentityStatus   IdentityStatus `json:"-"`
	GoogleSocialID   string         `json:"-"`
	FacebookSocialID string         `json:"-"`
	Email            string         `json:"-"`
	PasswordHash     string         `json:"-"`
	CreatedAt        time.Time      `json:"-"`
	UpdatedAt        time.Time      `json:"-"`
}

// SetEmail set email to User model
func (i *Identity) SetEmail(email string) error {
	fixedEmail := html.EscapeString(strings.TrimSpace(strings.ToLower(email)))

	if fixedEmail == "" {
		return ErrIdentityEmptyEmail
	}

	_, err := addressParser.Parse(fixedEmail)

	if err != nil {
		return ErrIdentityIncorrectEmail
	}

	i.Email = fixedEmail
	return nil
}

func (i *Identity) SetIdentityType(identityType IdentityType) {
	i.IdentityType = identityType
}

func (i *Identity) SetIdentityStatus(status IdentityStatus) {
	i.IdentityStatus = status
}

func CheckIdentityType(identityType string) bool {
	switch identityType {
	case IdentityTypeEmail.String():
		fallthrough
	case IdentityTypeFacebook.String():
		fallthrough
	case IdentityTypeGoogle.String():
		return true
	}

	return false
}

func CheckIdentityStatus(status string) bool {
	switch status {
	case IdentityStatusConfirmed.String():
		fallthrough
	case IdentityStatusUnconfirmed.String():
		return true
	}

	return false
}
