package manager

import (
	"database/sql"

	"github.com/artem-malko/auth-and-go/services/token"

	"github.com/artem-malko/auth-and-go/managers/user"
	"github.com/artem-malko/auth-and-go/models"
	"github.com/artem-malko/auth-and-go/services/account"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

func (m *userManager) DeleteUserByID(userID uuid.UUID) error {
	return m.runWithTransaction("DeleteUserByID", func(tx *sql.Tx) (err error) {
		err = m.sessionService.DeleteAllSessionsByAccountID(tx, userID)

		if err != nil {
			return errors.Wrap(err, "user manager: DeleteUserByID")
		}

		err = m.identityService.DeleteIdentitiesByAccountID(tx, userID)

		if err != nil {
			return errors.Wrap(err, "user manager: DeleteUserByID")
		}

		err = m.accountService.DeactivateAccountByID(userID, tx)

		if err != nil {
			switch errors.Cause(err) {
			case account.ErrNoAccountsUpdated:
				return user.ErrUserIsNotUpdated
			default:
				return errors.Wrap(err, "user manager: DeleteUserByID")
			}
		}

		return nil
	})
}

func (m *userManager) DeleteSessionBySessionID(sessionID uuid.UUID) error {
	err := m.sessionService.DeleteSessionBySessionID(m.db, sessionID)

	return errors.Wrap(err, "user manager: DeleteUsedTokens")
}

func (m *userManager) DeleteExpiredRegistrationConfirmations() error {
	return m.runWithTransaction("DeleteExpiredRegistrationConfirmations", func(tx *sql.Tx) error {
		expiredTokens, err := m.tokenService.DeleteExpiredTokens(models.TokenTypeRegistrationConfirmation, tx)

		if err != nil {
			switch errors.Cause(err) {
			case token.ErrNoTokensUpdated:
				return nil
			default:
				return errors.Wrap(err, "user manager: DeleteExpiredRegistrationConfirmations")
			}
		}

		var identityIDs, accountIDs []uuid.UUID

		for _, t := range expiredTokens {
			identityIDs = append(identityIDs, t.IdentityID)
			accountIDs = append(accountIDs, t.AccountID)
		}

		err = m.identityService.DeleteIdentitiesByIdentityIDs(identityIDs, tx)

		if err != nil {
			return errors.Wrap(err, "user manager: DeleteExpiredRegistrationConfirmations")
		}

		err = m.accountService.DeleteUnconfirmedAccountsByAccountIDs(accountIDs, tx)

		return errors.Wrap(err, "user manager: DeleteExpiredRegistrationConfirmations")
	})
}

func (m *userManager) DeleteUsedTokens() error {
	err := m.tokenService.DeleteUsedTokens()

	return errors.Wrap(err, "user manager: DeleteUsedTokens")
}

func (m *userManager) DeleteExpiredToken(tokenType models.TokenType) error {
	if tokenType == models.TokenTypeRegistrationConfirmation {
		return errors.New("Use DeleteExpiredRegistrationConfirmations for registration confirmation token deleting")
	}

	_, err := m.tokenService.DeleteExpiredTokens(tokenType, nil)

	if err != nil {
		switch errors.Cause(err) {
		case token.ErrNoTokensUpdated:
			return nil
		default:
			return errors.Wrap(err, "user manager: DeleteExpiredEmailConfirmations")
		}
	}

	return nil
}

func (m *userManager) DeleteExpiredSessions() error {
	err := m.sessionService.DeleteExpiredSessions()

	return errors.Wrap(err, "user manager: DeleteExpiredSessions")
}
