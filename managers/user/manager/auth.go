package manager

import (
	"database/sql"

	"github.com/artem-malko/auth-and-go/credentials"
	"github.com/artem-malko/auth-and-go/services/session"
	"github.com/artem-malko/auth-and-go/services/token"

	"github.com/artem-malko/auth-and-go/managers/user"
	"github.com/artem-malko/auth-and-go/models"
	"github.com/artem-malko/auth-and-go/services/identity"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

func (m *userManager) updateAccountAndGetSessionTokens(
	accountID, identityID uuid.UUID,
	clientIP string,
	clientID models.ClientID,
) (*models.Account, *models.SessionTokens, error) {
	var updatedAccount *models.Account
	var sessionTokens *models.SessionTokens

	transactionErr := m.runWithTransaction("updateAccountAndGetSessionTokens", func(tx *sql.Tx) (err error) {
		updatedAccount, err = m.accountService.LoginAccount(accountID, clientIP, tx)

		if err != nil {
			return errors.Wrap(err, "user manager: updateAccountAndGetSessionTokens error")
		}

		s, err := m.sessionService.CreateSession(accountID, identityID, clientID, tx)

		if err != nil {
			return errors.Wrap(err, "user manager: CreateSession error")
		}

		sessionTokens, err = credentials.CreateSessionTokens(*s, m.accessTokenSecretKey, m.refreshTokenSecretKey)

		if err != nil {
			return errors.Wrap(err, "user manager: createSessionTokens error")
		}

		return nil
	})

	if transactionErr != nil {
		return nil, nil, transactionErr
	}

	return updatedAccount, sessionTokens, nil
}

func (m *userManager) LoginWithEmailAndPassword(email, password, clientIP string, clientID models.ClientID) (*models.User, *models.SessionTokens, error) {
	i, err := m.identityService.GetEmailIdentityByEmailAndPassword(m.db, email, password)

	if err != nil {
		switch errors.Cause(err) {
		case identity.ErrNoIdentityFound:
			return nil, nil, user.ErrUserNoIdentitiesFound
		}
		return nil, nil, errors.Wrap(err, "userManager: GetEmailIdentityByEmailAndPassword error")
	}

	// Only confirmed identity can be proceeded
	if i.IdentityStatus != models.IdentityStatusConfirmed {
		return nil, nil, user.ErrUserNoIdentitiesFound
	}

	a, sessionTokens, err := m.updateAccountAndGetSessionTokens(i.AccountID, i.ID, clientIP, clientID)

	if err != nil {
		return nil, nil, errors.Wrap(err, "userManager: updateAccountAndGetSessionTokens error")
	}

	newUser := models.NewUser(models.UserTypeFull, *a, []*models.Identity{i})

	return &newUser, sessionTokens, nil
}

func (m *userManager) RefreshSession(refreshToken uuid.UUID) (*models.SessionTokens, error) {
	refreshedSession, err := m.sessionService.RefreshSession(refreshToken)

	if err != nil {
		switch errors.Cause(err) {
		case session.ErrNoSessionsUpdated:
			return nil, user.ErrUserSessionNotFound
		default:
			return nil, err
		}
	}

	sessionTokens, err := credentials.CreateSessionTokens(*refreshedSession, m.accessTokenSecretKey, m.refreshTokenSecretKey)

	if err != nil {
		return nil, errors.Wrap(err, "userManager: createSessionTokens in RefreshSession error")
	}

	return sessionTokens, nil
}

func (m *userManager) ConfirmRegistration(confirmationToken uuid.UUID) (*models.User, *models.SessionTokens, error) {
	var newUser models.User
	var newSessionTokens *models.SessionTokens

	transactionErr := m.runWithTransaction("ConfirmRegistration", func(tx *sql.Tx) (err error) {
		tokenInfo, err := m.tokenService.Use(tx, confirmationToken)

		if err != nil {
			switch errors.Cause(err) {
			case token.ErrNoTokensUpdated:
				return user.ErrUserIncorrectTokenToUse
			default:
				return errors.Wrap(err, "user Manager: ConfirmRegistration")
			}
		}

		a, err := m.accountService.ConfirmAccount(tokenInfo.AccountID, tx)

		// Every error from ConfirmAccount is fatal
		if err != nil {
			return errors.Wrap(err, "user Manager: ConfirmRegistration")
		}

		i, err := m.identityService.ConfirmIdentity(tx, tokenInfo.IdentityID)

		// Every error from ConfirmIdentity is fatal
		if err != nil {
			return errors.Wrap(err, "user Manager: ConfirmRegistration")
		}

		s, err := m.sessionService.CreateSession(tokenInfo.AccountID, tokenInfo.IdentityID, tokenInfo.ClientID, tx)

		// Every error from CreateSession is fatal
		if err != nil {
			return errors.Wrap(err, "user Manager: ConfirmRegistration")
		}

		newSessionTokens, err = credentials.CreateSessionTokens(*s, m.accessTokenSecretKey, m.refreshTokenSecretKey)

		// Every error from CreateSession is fatal
		if err != nil {
			return errors.Wrap(err, "user Manager: CreateSessionTokens")
		}

		sendEmailErr := m.mailerService.SendRegistrationConfirmedEmail(i.Email, a.Profile.FirstName, a.AccountName)

		// sendEmailErr should not rollback transaction
		// it is not required to send email on successful registration confirmation
		if sendEmailErr != nil {
			m.logger.
				WithField("source", "user_manager").
				Error(
					errors.Wrap(sendEmailErr, "user manager send registration confirmed email").Error(),
				)
		}

		newUser = models.NewUser(models.UserTypeFull, *a, []*models.Identity{i})

		return nil
	})

	if transactionErr != nil {
		return nil, nil, transactionErr
	}

	return &newUser, newSessionTokens, nil
}

func (m *userManager) ContinueWithOAuth(params user.ContinueWithOAuthParams) (*models.SessionTokens, error) {
	identityBySocialID, err := m.identityService.GetIdentityBySocialID(
		m.db,
		params.SocialID,
		params.SocialNetworkType,
	)

	if err != nil {
		switch errors.Cause(err) {
		case identity.ErrNoIdentityFound:
			return m.loginNewOAuthUser(params)
		default:
			return nil, errors.Wrap(err, "user manager: ContinueWithOAuth error")
		}
	}

	return m.loginExistedOAuthUser(params, identityBySocialID)
}

func (m *userManager) loginExistedOAuthUser(
	params user.ContinueWithOAuthParams,
	identity *models.Identity,
) (*models.SessionTokens, error) {
	tx, err := m.db.Begin()
	if err != nil {
		return nil, errors.Wrap(err, "user manager: loginExistedOAuthUser error during transaction opening")
	}

	defer func() {
		if err != nil {
			rollbackErr := tx.Rollback()

			if rollbackErr != nil {
				err = errors.Wrap(rollbackErr, "user manager: rollbackErr loginExistedOAuthUser")
			}
		}
	}()

	s, err := m.sessionService.CreateSession(identity.AccountID, identity.ID, params.ClientID, nil)

	// Every error from CreateSession is fatal
	if err != nil {
		return nil, errors.Wrap(err, "user manager: loginExistedOAuthUser CreateSession error")
	}

	sessionTokens, err := credentials.CreateSessionTokens(*s, m.accessTokenSecretKey, m.refreshTokenSecretKey)

	// Every error from CreateSessionTokens is fatal
	if err != nil {
		return nil, errors.Wrap(err, "user manager: loginExistedOAuthUser CreateSessionTokens error")
	}

	_, err = m.accountService.LoginAccount(identity.AccountID, params.ClientIP, tx)

	// Every error from LoginAccount is fatal
	if err != nil {
		return nil, errors.Wrap(err, "user manager: loginExistedOAuthUser LoginAccount error")
	}

	err = tx.Commit()

	if err != nil {
		return nil, errors.Wrap(err, "user manager: loginExistedOAuthUser transaction commit error")
	}

	return sessionTokens, nil
}

func (m *userManager) loginNewOAuthUser(params user.ContinueWithOAuthParams) (*models.SessionTokens, error) {
	tx, err := m.db.Begin()
	if err != nil {
		return nil, errors.Wrap(err, "user manager: loginNewOAuthUser error during transaction opening")
	}

	defer func() {
		if err != nil {
			rollbackErr := tx.Rollback()

			if rollbackErr != nil {
				err = errors.Wrap(rollbackErr, "user manager: rollbackErr loginNewOAuthUser")
			}
		}
	}()

	var identitiesByEmail []*models.Identity
	var errGetIdentityByEmail error
	accountID := uuid.New()

	if params.Email != "" {
		identitiesByEmail, errGetIdentityByEmail = m.identityService.GetIdentitiesByEmail(m.db, params.Email)

		// Only fatal error for request will be returned there
		if errGetIdentityByEmail != nil && errors.Cause(errGetIdentityByEmail) != identity.ErrNoIdentityFound {
			return nil, errors.Wrap(errGetIdentityByEmail, "user manager: loginNewOAuthUser error")
		}
	}

	// Try to find any identities with the email from current social login
	if len(identitiesByEmail) != 0 {
		accountID = identitiesByEmail[0].AccountID

		// We can confirm account, cause user from social network is not a robot, hopefully =)
		_, err = m.accountService.ConfirmAccount(accountID, tx)

		// Every error from ConfirmAccount is fatal
		if err != nil {
			return nil, errors.Wrap(err, "user manager: loginNewOAuthUser ConfirmAccount error")
		}
	} else {
		// If there is no identity for email from social network — create account
		accountByEmailFromSocialLogin, err := m.accountService.CreateAccountWithOAuth(
			params.Email, params.FirstName, params.LastName, params.AvatarURL,
			tx,
		)

		// Every error from CreateAccountWithOAuth is fatal
		if err != nil {
			return nil, errors.Wrap(err, "user manager: loginNewOAuthUser CreateAccountWithOAuth error")
		}

		accountID = accountByEmailFromSocialLogin.ID
	}

	socialIdentity, err := m.identityService.CreateOAuthIdentity(
		tx,
		accountID,
		params.SocialID,
		params.SocialNetworkType,
		params.Email,
	)

	// Every error from CreateOAuthIdentity is fatal
	if err != nil {
		return nil, errors.Wrap(err, "user manager: loginNewOAuthUser CreateOAuthIdentity error")
	}

	s, err := m.sessionService.CreateSession(tx, accountID, socialIdentity.ID, params.ClientID)

	// Every error from CreateSession is fatal
	if err != nil {
		return nil, errors.Wrap(err, "user manager: loginNewOAuthUser CreateSession error")
	}

	sessionTokens, err := credentials.CreateSessionTokens(*s, m.accessTokenSecretKey, m.refreshTokenSecretKey)

	// Every error from CreateSessionTokens is fatal
	if err != nil {
		return nil, errors.Wrap(err, "user manager: loginNewOAuthUser CreateSessionTokens error")
	}

	_, err = m.accountService.LoginAccount(tx, accountID, params.ClientIP)

	// Every error from LoginAccount is fatal
	if err != nil {
		return nil, errors.Wrap(err, "user manager: loginNewOAuthUser LoginAccount error")
	}

	err = tx.Commit()

	if err != nil {
		return nil, errors.Wrap(err, "user manager: loginNewOAuthUser transaction commit error")
	}

	return sessionTokens, nil
}
