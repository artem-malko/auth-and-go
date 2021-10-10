package service

import (
	"github.com/artem-malko/auth-and-go/infrastructure/database"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

func (s *identityService) DeleteIdentitiesByAccountID(
	executor database.QueryExecutor,
	accountID uuid.UUID,
) (err error) {
	err = s.identityRepository.DeleteIdentitiesByAccountID(executor, accountID)

	return errors.Wrap(err, "identity service: DeleteIdentitiesByAccountID")
}

func (s *identityService) DeleteIdentitiesByIdentityIDs(
	executor database.QueryExecutor,
	identityIDs []uuid.UUID,
) error {
	err := s.identityRepository.DeleteIdentitiesByIdentityIDs(executor, identityIDs)

	return errors.Wrap(err, "identity service: DeleteIdentitiesByIdentityIDs")
}
