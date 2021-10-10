package constants

import "github.com/google/uuid"

type ChallengersOfTheWorldUser struct {
	ID       uuid.UUID
	UserName string
	Email    string
}

type values struct {
	AccessTokenCookieName                        string
	AccessTokenMaxAgeInSeconds                   int
	RefreshTokenCookieName                       string
	RefreshTokenMaxAgeInSeconds                  int
	RegistrationConfirmationTokenMaxAgeInSeconds int
	ChangeEmailConfirmationTokenMaxAgeInSeconds  int
	AutoLoginTokenMaxAgeInSeconds                int
	DeleteUnconfirmedRegistrationTokensTimeout   int
	DeleteUnconfirmedEmailTokensTimeout          int
	DeleteAutoLoginTokensTimeout                 int
	DeleteUsedTokensTimeout                      int
	DeleteExpiredSessionsTimeout                 int
}

var Values = values{
	AccessTokenCookieName: "access_token",
	// 1 day
	AccessTokenMaxAgeInSeconds: 24 * 60 * 60,
	RefreshTokenCookieName:     "refresh_token",
	// 14 days
	RefreshTokenMaxAgeInSeconds: 14 * 24 * 60 * 60,
	// 1 hour
	RegistrationConfirmationTokenMaxAgeInSeconds: 60 * 60,
	// 31 min
	DeleteUnconfirmedRegistrationTokensTimeout: 31 * 60,
	// 1 day
	ChangeEmailConfirmationTokenMaxAgeInSeconds: 24 * 60 * 60,
	// 11 hours
	DeleteUnconfirmedEmailTokensTimeout: 11 * 60 * 60,
	// 14 days
	AutoLoginTokenMaxAgeInSeconds: 14 * 24 * 60 * 60,
	// 23 hours
	DeleteAutoLoginTokensTimeout: 23 * 60 * 60,
	// 24 hours
	DeleteUsedTokensTimeout: 6 * 60 * 60,
	// 14 min
	DeleteExpiredSessionsTimeout: 14 * 60,
}
