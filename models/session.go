package models

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrSessionIncorrectClientID = errors.New("client ID is incorrect")
)

type ClientID string

func (c ClientID) String() string {
	return string(c)
}

const (
	ClientIDWEB           ClientID = "web"
	ClientIDNativeAndroid ClientID = "native_android"
	ClientIDNativeIOS     ClientID = "native_ios"
)

type Session struct {
	ID                      uuid.UUID
	AccountID               uuid.UUID
	IdentityID              uuid.UUID
	ClientID                ClientID
	AccessToken             uuid.UUID
	AccessTokenExpiresDate  time.Time
	RefreshToken            uuid.UUID
	RefreshTokenExpiresDate time.Time
}

type SessionTokens struct {
	AccessToken  string
	RefreshToken string
}

func CheckClientID(clientID string) bool {
	switch clientID {
	case ClientIDWEB.String():
		fallthrough
	case ClientIDNativeAndroid.String():
		fallthrough
	case ClientIDNativeIOS.String():
		return true
	}

	return false
}
