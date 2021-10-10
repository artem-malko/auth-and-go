package service

import (
	"time"

	"github.com/artem-malko/auth-and-go/infrastructure/database"
	"github.com/artem-malko/auth-and-go/models"
	identityRepository "github.com/artem-malko/auth-and-go/repositories/identity"
	"github.com/artem-malko/auth-and-go/services/identity"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

func (s *identityService) CreateEmailIdentity(
	executor database.QueryExecutor,
	accountID uuid.UUID,
	email, password string,
	identityStatus models.IdentityStatus,
) (*models.Identity, error) {
	newIdentity := models.Identity{
		ID:             uuid.New(),
		AccountID:      accountID,
		IdentityType:   models.IdentityTypeEmail,
		IdentityStatus: identityStatus,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	err := newIdentity.SetEmail(email)

	if err != nil {
		return nil, models.ErrIdentityIncorrectEmail
	}

	passwordHash, err := createHash(password)

	if err != nil {
		return nil, errors.Wrap(identity.ErrIdentityPasswordGeneration, err.Error())
	}

	newIdentity.PasswordHash = passwordHash

	createdIdentity, err := s.identityRepository.CreateEmailIdentity(executor, newIdentity)

	if err != nil {
		switch errors.Cause(err) {
		case identityRepository.ErrRepositoryEmailConstraint:
			return nil, identity.ErrIdentityExists
		default:
			return nil, err
		}

	}

	return createdIdentity, nil
}

func (s *identityService) CreateOAuthIdentity(
	executor database.QueryExecutor,
	accountID uuid.UUID,
	socialID string,
	socialNetworkType models.SocialNetworkType,
	email string,
) (*models.Identity, error) {
	newIdentity := models.Identity{
		ID:             uuid.New(),
		AccountID:      accountID,
		IdentityStatus: models.IdentityStatusConfirmed,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if socialNetworkType == models.SocialNetworkTypeGoogle {
		newIdentity.GoogleSocialID = socialID
		newIdentity.IdentityType = models.IdentityTypeGoogle
	}

	if socialNetworkType == models.SocialNetworkTypeFacebook {
		newIdentity.FacebookSocialID = socialID
		newIdentity.IdentityType = models.IdentityTypeFacebook
	}

	err := newIdentity.SetEmail(email)

	if err != nil {
		return nil, models.ErrIdentityIncorrectEmail
	}

	createdIdentity, err := s.identityRepository.CreateSocialIdentity(executor, newIdentity)

	if err != nil {
		switch errors.Cause(err) {
		case identityRepository.ErrRepositorySocialFacebookConstraint:
			fallthrough
		case identityRepository.ErrRepositorySocialGoogleConstraint:
			return nil, identity.ErrIdentityExists
		default:
			return nil, errors.Wrap(err, "identity service: CreateOAuthIdentity")
		}

	}

	return createdIdentity, nil
}
