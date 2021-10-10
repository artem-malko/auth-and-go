package manager

import (
	"database/sql"

	"github.com/artem-malko/auth-and-go/managers/user"
	"github.com/artem-malko/auth-and-go/models"
	"github.com/artem-malko/auth-and-go/services/identity"
	"github.com/pkg/errors"
)

func (m *userManager) CreateUserWithEmail(email, password string, clientID models.ClientID) error {
	return m.runWithTransaction("CreateUserWithEmail", func(tx *sql.Tx) error {
		var account *models.Account
		var existedIdentityWithEmailFromRequest *models.Identity

		identitiesByEmail, err := m.identityService.GetIdentitiesByEmail(tx, email)

		if err != nil {
			switch errors.Cause(err) {
			case identity.ErrNoIdentityFound:
				// It's ok, so we need to create that identity
			default:
				return errors.Wrap(err, "user manager: CreateUserWithEmail GetEmailIdentityByEmail error")
			}
		}

		// If there is an identity with the same email, just use its account
		if len(identitiesByEmail) != 0 {
			// We need to filter users, who want to create account by using already registered email
			for _, i := range identitiesByEmail {
				if i.IdentityType == models.IdentityTypeEmail {
					if i.IdentityStatus == models.IdentityStatusConfirmed {
						// @TODO correct name
						err = m.mailerService.SendUnexpectedRegistrationEmail(email, "friend", "")

						if err != nil {
							return errors.Wrap(err, "user manager: CreateUserWithEmail error")
						}
					}

					if i.IdentityStatus == models.IdentityStatusUnconfirmed {
						activeToken, err := m.tokenService.GetActiveTokenByIdentityID(
							tx, i.ID, models.TokenTypeRegistrationConfirmation,
						)

						if err != nil {
							return errors.Wrap(err, "user manager: CreateUserWithEmail error")
						}

						err = m.mailerService.SendRegistrationConfirmationEmail(email, activeToken.ID)

						if err != nil {
							return errors.Wrap(err, "user manager: CreateUserWithEmail error")
						}
					}

					return user.ErrUserWithSameIdentityExists
				}

				// Other identities than email can be used as existed
				existedIdentityWithEmailFromRequest = i
			}
		}

		if existedIdentityWithEmailFromRequest != nil {
			// Ok, there is no any email identity, we can get accountID from any other
			// We can choose any identity, cause all of them are from one person, hopefully =)
			account, err = m.accountService.GetAccountByID(tx, existedIdentityWithEmailFromRequest.AccountID)

			// It can not be possible to not have account, if there is an identity for it
			// But handle this too
			if err != nil {
				return errors.Wrap(err, "user manager: CreateUserWithEmail GetAccountByID error")
			}
		} else {
			account, err = m.accountService.CreateAccount(tx, email, models.AccountStatusUnconfirmed)

			if err != nil {
				return errors.Wrap(err, "user manager: CreateUserWithEmail error")
			}
		}

		newIdentity, err := m.identityService.CreateEmailIdentity(tx, account.ID, email, password, models.IdentityStatusUnconfirmed)

		if err != nil {
			switch errors.Cause(err) {
			case identity.ErrIdentityExists:
				err = m.mailerService.SendUnexpectedRegistrationEmail(email, account.Profile.FirstName, account.AccountName)

				if err != nil {
					return errors.Wrap(err, "user manager: CreateUserWithEmail error")
				}

				return user.ErrUserWithSameIdentityExists
			default:
				return errors.Wrap(err, "user manager: CreateUserWithEmail")
			}
		}

		token, err := m.tokenService.Create(
			tx,
			models.TokenTypeRegistrationConfirmation,
			clientID,
			account.ID,
			newIdentity.ID,
		)

		if err != nil {
			return errors.Wrap(err, "user manager: CreateUserWithEmail")
		}

		err = m.mailerService.SendRegistrationConfirmationEmail(email, token.ID)

		if err != nil {
			return errors.Wrap(err, "user manager: CreateUserWithEmail")
		}

		return nil
	})
}
