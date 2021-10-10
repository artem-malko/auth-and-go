package service

import (
	"github.com/artem-malko/auth-and-go/infrastructure/database"
	"github.com/artem-malko/auth-and-go/models"
	"github.com/artem-malko/auth-and-go/services/identity"
	"github.com/pkg/errors"
)

func (s *identityService) GetIdentitiesByEmail(
	executor database.QueryExecutor,
	email string,
) ([]*models.Identity, error) {
	identities, err := s.identityRepository.GetIdentitiesByEmail(executor, email)

	if err != nil {
		switch errors.Cause(err) {
		case database.ErrRepositoryNoEntitiesFound:
			return nil, identity.ErrNoIdentityFound
		default:
			return nil, errors.Wrap(err, "identity service: GetIdentitiesByEmail")
		}
	}

	return identities, nil
}

func (s *identityService) GetEmailIdentityByEmail(
	executor database.QueryExecutor,
	email string,
) (*models.Identity, error) {
	identityByEmail, err := s.identityRepository.GetEmailIdentityByEmail(executor, email)

	if err != nil {
		switch errors.Cause(err) {
		case database.ErrRepositoryNoEntitiesFound:
			return nil, identity.ErrNoIdentityFound
		default:
			return nil, errors.Wrap(err, "identity service: GetEmailIdentityByEmail")
		}
	}

	return identityByEmail, nil
}

func (s *identityService) GetEmailIdentityByEmailAndPassword(
	executor database.QueryExecutor,
	email, password string,
) (*models.Identity, error) {
	identityByEmail, err := s.identityRepository.GetEmailIdentityByEmail(executor, email)

	if err != nil {
		switch errors.Cause(err) {
		case database.ErrRepositoryNoEntitiesFound:
			return nil, identity.ErrNoIdentityFound
		default:
			return nil, errors.Wrap(err, "identity service: GetEmailIdentityByEmailAndPassword")
		}
	}

	if !validatePassword(password, identityByEmail.PasswordHash) {
		return nil, identity.ErrNoIdentityFound
	}

	return identityByEmail, nil
}

func (s *identityService) GetIdentityBySocialID(
	executor database.QueryExecutor,
	socialID string,
	socialNetworkType models.SocialNetworkType,
) (*models.Identity, error) {
	identityBySocialID, err := s.identityRepository.GetIdentityBySocialID(
		executor,
		socialID,
		socialNetworkType,
	)

	if err != nil {
		switch errors.Cause(err) {
		case database.ErrRepositoryNoEntitiesFound:
			return nil, identity.ErrNoIdentityFound
		default:
			return nil, errors.Wrap(err, "identity service: GetIdentityBySocialID")
		}
	}

	return identityBySocialID, nil
}
