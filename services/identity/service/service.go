package service

import (
	"github.com/artem-malko/auth-and-go/infrastructure/database"
	"github.com/artem-malko/auth-and-go/services/identity"
	"github.com/pkg/errors"

	"github.com/artem-malko/auth-and-go/models"
	"github.com/google/uuid"
)

type identityService struct {
	identityRepository identity.Repository
}

func New(identityRepository identity.Repository) identity.Service {
	return &identityService{
		identityRepository: identityRepository,
	}
}

func (s *identityService) GetIdentitiesByAccountID(
	executor database.QueryExecutor,
	accountID uuid.UUID,
) ([]*models.Identity, error) {
	identities, err := s.identityRepository.GetIdentitiesByAccountID(
		executor,
		accountID,
	)

	if err != nil {
		return nil, errors.Wrap(err, "identity service: GetIdentitiesByAccountID")
	}

	return identities, nil
}
