package models

import (
	"encoding/json"
	"time"

	"github.com/artem-malko/auth-and-go/infrastructure/convert"

	"github.com/google/uuid"
)

type UserType string

const (
	UserTypeFull  UserType = "full"
	UserTypeShort UserType = "short"
)

func convertProfileToMap(profile Profile, avatarPicture File, mediaPrefix string) map[string]interface{} {
	return map[string]interface{}{
		"first_name":   convert.NewStringPointer(profile.FirstName),
		"last_name":    convert.NewStringPointer(profile.LastName),
		"gender":       convert.NewStringPointer(profile.Gender),
		"country":      convert.NewStringPointer(profile.Country),
		"city":         convert.NewStringPointer(profile.City),
		"description":  convert.NewStringPointer(profile.Description),
		"avatar":       avatarPicture.ConvertToMap(mediaPrefix),
		"birthday":     convert.NewTimePointer(profile.Birthday),
		"social_links": profile.SocialLinks,
		"interests":    profile.Interests,
	}
}

func convertSettingsToMap(settings Settings) map[string]interface{} {
	return map[string]interface{}{
		"notifications": map[string]interface{}{
			"email": map[string]interface{}{
				"email": settings.Notifications.Email.Email,
				"scope": settings.Notifications.Email.Scope,
			},
		},
	}
}

type User struct {
	ID            uuid.UUID
	AccountType   AccountType
	AccountStatus AccountStatus
	AccountName   string
	Profile       Profile
	Identities    []UserIdentity
	LastIP        string
	LastLogin     time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Settings      Settings
	UserType      UserType
	AvatarFile    File
	MediaPrefix   string
}

type UserIdentity struct {
	IdentityType   IdentityType
	Email          string
	IdentityStatus IdentityStatus
}

func (u *User) MarshalJSON() ([]byte, error) {
	var identities []map[string]interface{}

	for _, identity := range u.Identities {
		identities = append(identities, map[string]interface{}{
			"type":   identity.IdentityType,
			"status": identity.IdentityStatus,
			"email":  convert.NewStringPointer(identity.Email),
		})
	}

	userToMarshal := map[string]interface{}{
		"id":      u.ID,
		"type":    u.AccountType,
		"name":    u.AccountName,
		"profile": convertProfileToMap(u.Profile, u.AvatarFile, u.MediaPrefix),
	}

	settingsToMarshal := convertSettingsToMap(u.Settings)

	if u.UserType == UserTypeFull {
		userToMarshal["status"] = u.AccountStatus
		userToMarshal["identities"] = identities
		userToMarshal["last_login"] = u.LastLogin.Format(time.RFC3339)
		userToMarshal["last_ip"] = u.LastIP
		userToMarshal["created_at"] = u.CreatedAt.Format(time.RFC3339)
		userToMarshal["updated_at"] = u.UpdatedAt.Format(time.RFC3339)
		userToMarshal["settings"] = settingsToMarshal
	}

	return json.Marshal(userToMarshal)
}

func NewUser(
	userType UserType,
	account Account,
	identities []*Identity,
) User {
	var userIdentities []UserIdentity

	for _, identity := range identities {
		userIdentities = append(userIdentities, UserIdentity{
			IdentityType:   identity.IdentityType,
			Email:          identity.Email,
			IdentityStatus: identity.IdentityStatus,
		})
	}

	newUser := User{
		ID:            account.ID,
		AccountType:   account.AccountType,
		AccountStatus: account.AccountStatus,
		AccountName:   account.AccountName,
		Profile:       account.Profile,
		Settings:      account.Settings,
		LastIP:        account.LastIP,
		LastLogin:     account.LastLogin,
		CreatedAt:     account.CreatedAt,
		UpdatedAt:     account.UpdatedAt,
		Identities:    userIdentities,
		UserType:      userType,
		AvatarURL:     account.AvatarURL,
	}

	return newUser
}
