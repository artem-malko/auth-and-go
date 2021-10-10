package manager

import (
	"github.com/artem-malko/auth-and-go/managers/user"
	"github.com/artem-malko/auth-and-go/models"
	"github.com/artem-malko/auth-and-go/services/account"
	"github.com/artem-malko/auth-and-go/services/identity"
	"github.com/artem-malko/auth-and-go/services/session"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

// GetUserByID retrieve account info and its identities
func (m *userManager) GetUserByID(userID uuid.UUID) (*models.User, error) {
	a, err := m.accountService.GetAccountByID(userID)

	if err != nil {
		switch errors.Cause(err) {
		case account.ErrAccountNotFound:
			return nil, user.ErrUserNotFound
		default:
			return nil, errors.Wrap(err, "userManager: can't get info for account with ID "+userID.String())
		}
	}

	u := models.NewUser(models.UserTypeShort, *a, nil)

	return &u, nil
}

func (m *userManager) GetUserByName(userName string) (*models.User, error) {
	a, err := m.accountService.GetAccountByName(userName)

	if err != nil {
		switch errors.Cause(err) {
		case account.ErrAccountNotFound:
			return nil, user.ErrUserNotFound
		default:
			return nil, errors.Wrap(err, "userManager: can't get info for account with Name "+userName)
		}
	}

	u := models.NewUser(models.UserTypeShort, *a, nil)

	return &u, nil
}

func (m *userManager) GetSessionByAccessToken(accessToken uuid.UUID) (*models.Session, error) {
	sessionByAccessToken, err := m.sessionService.GetSessionByAccessToken(accessToken)

	if err != nil {
		switch errors.Cause(err) {
		case session.ErrSessionNotFound:
			return nil, user.ErrUserSessionNotFound
		default:
			return nil, errors.Wrap(err, "user manager: GetSessionByAccessToken err")
		}
	}

	return sessionByAccessToken, nil
}

func (m *userManager) GetFullUser(userID uuid.UUID) (*models.User, error) {
	a, err := m.accountService.GetAccountByID(userID)

	if err != nil {
		switch errors.Cause(err) {
		case account.ErrAccountNotFound:
			return nil, user.ErrUserNotFound
		default:
			return nil, errors.Wrap(err, "userManager: GetAccountByID in GetFullUser error")
		}
	}

	i, err := m.identityService.GetIdentitiesByAccountID(userID)

	if err != nil {
		switch errors.Cause(err) {
		case identity.ErrNoIdentityFound:
			return nil, user.ErrUserNotFound
		default:
			return nil, errors.Wrap(err, "userManager: GetIdentitiesByAccountID in GetFullUser error")
		}
	}

	u := models.NewUser(models.UserTypeFull, *a, i)

	return &u, nil
}
