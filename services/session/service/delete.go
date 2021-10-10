package service

import (
	"github.com/artem-malko/auth-and-go/infrastructure/database"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

func (s *sessionService) DeleteSessionBySessionID(
	executor database.QueryExecutor,
	sessionID uuid.UUID,
) error {
	err := s.sessionRepository.DeleteSessionBySessionID(executor, sessionID)

	return errors.Wrap(err, "session service: DeleteSessionBySessionID error")
}

func (s *sessionService) DeleteAllSessionsByAccountID(
	executor database.QueryExecutor,
	accountID uuid.UUID,
) error {
	err := s.sessionRepository.DeleteAllSessionsByAccountID(executor, accountID)

	return errors.Wrap(err, "session service: DeleteAllSessionsByAccountID error")
}

func (s *sessionService) DeleteExpiredSessions(executor database.QueryExecutor) error {
	err := s.sessionRepository.DeleteExpiredSessions(executor)

	return errors.Wrap(err, "session service: DeleteExpiredSessions error")
}
