package identity

import (
	"errors"

	"github.com/artem-malko/auth-and-go/infrastructure/database"

	"github.com/artem-malko/auth-and-go/models"
	"github.com/google/uuid"
)

var (
	ErrIdentityExists                     = errors.New("identity already exists")
	ErrRepositoryEmailConstraint          = errors.New("identity with the same email already exists")
	ErrRepositorySocialGoogleConstraint   = errors.New("google identity with the same ID already exists")
	ErrRepositorySocialFacebookConstraint = errors.New("facebook identity with the same ID already exists")
	ErrRepositoryUnknownSocialNetworkType = errors.New("unknown social network type")
)

// IdentityRepository is an interface for any Repository
type Repository interface {
	GetIdentitiesByAccountID(
		executor database.QueryExecutor,
		accountID uuid.UUID,
	) ([]*models.Identity, error)
	UpdateIdentityStatus(
		executor database.QueryExecutor,
		identityID uuid.UUID,
		identityStatus models.IdentityStatus,
	) (*models.Identity, error)
	CreateEmailIdentity(
		executor database.QueryExecutor,
		identity models.Identity,
	) (*models.Identity, error)
	CreateSocialIdentity(
		executor database.QueryExecutor,
		identity models.Identity,
	) (*models.Identity, error)
	GetIdentitiesByEmail(executor database.QueryExecutor, email string) ([]*models.Identity, error)
	GetEmailIdentityByEmail(executor database.QueryExecutor, email string) (*models.Identity, error)
	GetIdentityBySocialID(
		executor database.QueryExecutor,
		socialID string,
		socialNetworkType models.SocialNetworkType,
	) (*models.Identity, error)
	//CreateIdentity(user models.DBIdentity) (*models.Identity, error)
	UpdatePasswordHash(executor database.QueryExecutor, accountID uuid.UUID, passwordHash string) error
	DeleteIdentitiesByAccountID(executor database.QueryExecutor, accountID uuid.UUID) error
	DeleteIdentitiesByIdentityIDs(executor database.QueryExecutor, identityIDs []uuid.UUID) error
}
