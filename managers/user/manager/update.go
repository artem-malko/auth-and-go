package manager

import (
	"github.com/artem-malko/auth-and-go/managers/user"
	"github.com/artem-malko/auth-and-go/services/account"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

func (m *userManager) UpdateAccountName(userID uuid.UUID, name string) error {
	err := m.accountService.UpdateAccountName(m.db, userID, name)

	if err != nil {
		switch errors.Cause(err) {
		case account.ErrAccountNameExists:
			return user.ErrUserWithSameNameExists
		case account.ErrNoAccountsUpdated:
			return user.ErrUserIsNotUpdated
		default:
			return errors.Wrap(err, "user manager: UpdateAccountName")
		}
	}

	return nil
}
