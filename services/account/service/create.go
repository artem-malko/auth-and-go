package service

import (
	"time"

	"github.com/pkg/errors"

	"github.com/artem-malko/auth-and-go/infrastructure/database"
	"github.com/artem-malko/auth-and-go/services/account/service/generator"

	"github.com/artem-malko/auth-and-go/models"
	"github.com/google/uuid"
)

func (s *accountService) CreateAccount(
	executor database.QueryExecutor,
	email string,
	accountStatus models.AccountStatus,
) (*models.Account, error) {
	accountID := uuid.New()
	scope := []models.EmailNotificationScopeType{
		models.EmailNotificationScopeTypeBase,
	}

	if accountStatus == models.AccountStatusConfirmed {
		scope = append(
			scope,
			models.EmailNotificationScopeTypeNews,
			models.EmailNotificationScopeTypeReminders,
		)
	}

	newAccount := models.Account{
		ID:            accountID,
		AccountType:   models.AccountTypeFree,
		AccountStatus: accountStatus,
		AccountName:   generator.GenerateAccountName(accountID),
		LastLogin:     time.Time{},
		Profile: models.Profile{
			SocialLinks: []models.SocialLink{},
			Interests:   []string{},
		},
		Settings: models.Settings{
			Notifications: models.Notifications{
				Email: models.EmailNotifications{
					Email: email,
					Scope: scope,
				},
			},
		},
	}

	createdAccount, err := s.accountRepository.CreateAccount(executor, newAccount)

	if err != nil {
		return nil, errors.Wrap(err, "account service: CreateAccount")
	}

	return createdAccount, nil
}

func (s *accountService) CreateAccountWithOAuth(
	executor database.QueryExecutor,
	email, firstName, lastName, avatarURL string,
) (*models.Account, error) {
	accountID := uuid.New()

	// @TODO fix problem with empty email

	newAccount := models.Account{
		ID:            accountID,
		AccountType:   models.AccountTypeFree,
		AccountStatus: models.AccountStatusConfirmed,
		AccountName:   generator.GenerateAccountName(accountID),
		LastLogin:     time.Time{},
		Profile: models.Profile{
			FirstName:   firstName,
			LastName:    lastName,
			AvatarURL:   avatarURL,
			SocialLinks: []models.SocialLink{},
			Interests:   []string{},
		},
		Settings: models.Settings{
			Notifications: models.Notifications{
				Email: models.EmailNotifications{
					Email: email,
					Scope: []models.EmailNotificationScopeType{
						models.EmailNotificationScopeTypeBase,
						models.EmailNotificationScopeTypeNews,
						models.EmailNotificationScopeTypeReminders,
					},
				},
			},
		},
	}

	createdAccount, err := s.accountRepository.CreateAccount(executor, newAccount)

	if err != nil {
		return nil, errors.Wrap(err, "account service: CreateAccountWithOAuth")
	}

	return createdAccount, nil
}
