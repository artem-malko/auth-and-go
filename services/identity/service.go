package identity

import (
	"errors"

	"github.com/artem-malko/auth-and-go/infrastructure/database"
	"github.com/artem-malko/auth-and-go/models"
	"github.com/google/uuid"
)

var (
	ErrIdentityExists             = errors.New("identity exists")
	ErrIdentityPasswordGeneration = errors.New("password generation error")
	ErrNoIdentityFound            = errors.New("no identities found")
	ErrNoIdentityUpdated          = errors.New("identities with sent params are not updated")
)

type Repository interface {
	GetIdentitiesByAccountID(executor database.QueryExecutor, accountID uuid.UUID) ([]*models.Identity, error)
	UpdateIdentityStatus(executor database.QueryExecutor, identityID uuid.UUID, identityStatus models.IdentityStatus) (*models.Identity, error)
	CreateEmailIdentity(executor database.QueryExecutor, identity models.Identity) (*models.Identity, error)
	CreateSocialIdentity(executor database.QueryExecutor, identity models.Identity) (*models.Identity, error)
	GetIdentitiesByEmail(executor database.QueryExecutor, email string) ([]*models.Identity, error)
	GetEmailIdentityByEmail(executor database.QueryExecutor, email string) (*models.Identity, error)
	GetIdentityBySocialID(executor database.QueryExecutor, socialID string, socialNetworkType models.SocialNetworkType) (*models.Identity, error)
	//CreateIdentity(user models.DBIdentity) (*models.Identity, error)
	UpdatePasswordHash(executor database.QueryExecutor, accountID uuid.UUID, passwordHash string) error
	DeleteIdentitiesByAccountID(executor database.QueryExecutor, accountID uuid.UUID) error
	DeleteIdentitiesByIdentityIDs(executor database.QueryExecutor, identityIDs []uuid.UUID) error
}

type Service interface {
	GetIdentitiesByAccountID(executor database.QueryExecutor, accountID uuid.UUID) ([]*models.Identity, error)
	GetIdentityBySocialID(executor database.QueryExecutor, socialID string, socialNetworkType models.SocialNetworkType) (*models.Identity, error)
	GetIdentitiesByEmail(executor database.QueryExecutor, email string) ([]*models.Identity, error)
	GetEmailIdentityByEmail(executor database.QueryExecutor, email string) (*models.Identity, error)
	GetEmailIdentityByEmailAndPassword(executor database.QueryExecutor, email, password string) (*models.Identity, error)
	CreateEmailIdentity(
		executor database.QueryExecutor,
		accountID uuid.UUID,
		email, password string,
		identityStatus models.IdentityStatus,
	) (*models.Identity, error)
	CreateOAuthIdentity(
		executor database.QueryExecutor,
		accountID uuid.UUID,
		socialID string,
		socialNetworkType models.SocialNetworkType,
		email string,
	) (*models.Identity, error)
	DeleteIdentitiesByAccountID(executor database.QueryExecutor, accountID uuid.UUID) error
	DeleteIdentitiesByIdentityIDs(executor database.QueryExecutor, identityIDs []uuid.UUID) error
	ConfirmIdentity(executor database.QueryExecutor, identityID uuid.UUID) (*models.Identity, error)
}
