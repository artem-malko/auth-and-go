package user

import (
	"errors"

	"github.com/artem-malko/auth-and-go/models"
	"github.com/google/uuid"
)

var (
	ErrUserWithSameIdentityExists            = errors.New("a user with the same identity already exists")
	ErrUserWithSameUnconfirmedIdentityExists = errors.New("a user with the same unconfirmed identity already exists")
	ErrUserNotFound                          = errors.New("user not found")
	ErrUserIsNotUpdated                      = errors.New("user is not updated")
	ErrUserWithSameNameExists                = errors.New("a user with the same name already exists")
	ErrUserNoIdentitiesFound                 = errors.New("there is no identities found for passed params")
	ErrUserSessionNotFound                   = errors.New("there is no sessions found for passed params")
	ErrUserNoExpiredRegistraions             = errors.New("there is no expired registration found")
	ErrUserIncorrectTokenToUse               = errors.New("token can not be used")
)

type ContinueWithOAuthParams struct {
	SocialNetworkType models.SocialNetworkType
	SocialID          string
	Email             string
	IsEmailVerified   bool
	FirstName         string
	LastName          string
	AvatarURL         string
	ClientID          models.ClientID
	ClientIP          string
}

type Manager interface {
	// Auth
	CreateUserWithEmail(email, password string, clientID models.ClientID) error
	ConfirmRegistration(confirmationToken uuid.UUID) (*models.User, *models.SessionTokens, error)
	DeleteUserByID(userID uuid.UUID) error
	LoginWithEmailAndPassword(
		email, password, clientIP string,
		clientID models.ClientID,
	) (*models.User, *models.SessionTokens, error)
	ContinueWithOAuth(params ContinueWithOAuthParams) (*models.SessionTokens, error)
	RefreshSession(refreshToken uuid.UUID) (*models.SessionTokens, error)
	DeleteSessionBySessionID(sessionID uuid.UUID) error
	GetSessionByAccessToken(accessToken uuid.UUID) (*models.Session, error)

	// Used in cron jobs
	DeleteExpiredRegistrationConfirmations() error
	DeleteUsedTokens() error
	DeleteExpiredToken(tokenType models.TokenType) error
	DeleteExpiredSessions() error

	// Users get
	GetUserByID(userID uuid.UUID) (*models.User, error)
	GetUserByName(userName string) (*models.User, error)
	GetFullUser(userID uuid.UUID) (*models.User, error)

	// Own user updates
	UpdateAccountName(userID uuid.UUID, name string) error
}
