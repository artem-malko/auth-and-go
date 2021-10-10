package models

import (
	"time"

	"github.com/google/uuid"
)

type AccountStatus string

func (a AccountStatus) String() string {
	return string(a)
}

const (
	AccountStatusUnconfirmed AccountStatus = "unconfirmed"
	AccountStatusConfirmed   AccountStatus = "confirmed"
	AccountStatusDeleted     AccountStatus = "deleted"
	AccountStatusBanned      AccountStatus = "banned"
)

type AccountType string

func (a AccountType) String() string {
	return string(a)
}

const (
	AccountTypeFree       AccountType = "free"
	AccountTypeCommercial AccountType = "commercial"
)

type SocialLinkType string

func (s SocialLinkType) String() string {
	return string(s)
}

const (
	SocialLinkTypeFacebook  SocialLinkType = "facebook"
	SocialLinkTypeGoogle    SocialLinkType = "google"
	SocialLinkTypeInstagram SocialLinkType = "instagram"
	SocialLinkTypeTwitter   SocialLinkType = "twitter"
)

type SocialLink struct {
	Type SocialLinkType
	URL  string
}

type Profile struct {
	FirstName   string       `json:"-"`
	LastName    string       `json:"-"`
	Country     string       `json:"-"`
	City        string       `json:"-"`
	Gender      string       `json:"-"`
	Birthday    time.Time    `json:"-"`
	Description string       `json:"-"`
	SocialLinks []SocialLink `json:"-"`
	AvatarURL   string       `json:"-"`
	Interests   []string     `json:"-"`
}

type EmailNotificationScopeType string

func (e EmailNotificationScopeType) String() string {
	return string(e)
}

const (
	// Auth confirm, change email and so on
	EmailNotificationScopeTypeBase EmailNotificationScopeType = "base"
	// Reminders about started challenges
	EmailNotificationScopeTypeReminders EmailNotificationScopeType = "reminders"
	// News about new challenges, new achievements
	EmailNotificationScopeTypeNews EmailNotificationScopeType = "news"
)

type EmailNotifications struct {
	Email string                       `json:"-"`
	Scope []EmailNotificationScopeType `json:"-"`
}

type Notifications struct {
	Email EmailNotifications `json:"-"`
}

type Settings struct {
	Notifications Notifications `json:"-"`
}

type Account struct {
	ID            uuid.UUID     `json:"-"`
	AccountStatus AccountStatus `json:"-"`
	AccountType   AccountType   `json:"-"`
	AccountName   string        `json:"-"`
	Profile       Profile       `json:"-"`
	Settings      Settings      `json:"-"`
	LastIP        string        `json:"-"`
	LastLogin     time.Time     `json:"-"`
	CreatedAt     time.Time     `json:"-"`
	UpdatedAt     time.Time     `json:"-"`
}

func (a *Account) SetAccountStatus(accountStatus AccountStatus) {
	a.AccountStatus = accountStatus
}

func (a *Account) SetAccountType(accountType AccountType) {
	a.AccountType = accountType
}

func CheckAccountStatus(accountStatus string) bool {
	switch accountStatus {
	case AccountStatusUnconfirmed.String():
		fallthrough
	case AccountStatusConfirmed.String():
		fallthrough
	case AccountStatusDeleted.String():
		fallthrough
	case AccountStatusBanned.String():
		return true
	}

	return false
}

func CheckAccountType(accountType string) bool {
	switch accountType {
	case AccountTypeFree.String():
		fallthrough
	case AccountTypeCommercial.String():
		return true
	}

	return false
}
