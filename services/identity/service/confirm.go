package service

import (
	"github.com/artem-malko/auth-and-go/infrastructure/database"
	"github.com/artem-malko/auth-and-go/services/identity"
	"github.com/pkg/errors"

	"github.com/artem-malko/auth-and-go/models"
	"github.com/google/uuid"
)

func (s *identityService) ConfirmIdentity(
	executor database.QueryExecutor,
	identityID uuid.UUID,
) (*models.Identity, error) {
	confirmedIdentity, err := s.identityRepository.UpdateIdentityStatus(
		executor,
		identityID,
		models.IdentityStatusConfirmed,
	)

	if err != nil {
		switch errors.Cause(err) {
		case database.ErrRepositoryNoRowsAffected:
			return nil, identity.ErrNoIdentityUpdated
		default:
			return nil, errors.Wrap(err, "identity service: ConfirmIdentity error")
		}
	}

	return confirmedIdentity, nil
}
