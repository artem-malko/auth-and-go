package service

import (
	"github.com/artem-malko/auth-and-go/infrastructure/database"
	"github.com/artem-malko/auth-and-go/models"
	"github.com/artem-malko/auth-and-go/services/session"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

func (s *sessionService) GetSessionByAccessToken(
	executor database.QueryExecutor,
	accessToken uuid.UUID,
) (*models.Session, error) {
	sessionByAccessToken, err := s.sessionRepository.GetSessionByAccessToken(executor, accessToken)

	if err != nil {
		switch errors.Cause(err) {
		case database.ErrRepositoryNoEntitiesFound:
			return nil, session.ErrSessionNotFound
		default:
			return nil, errors.Wrap(err, "session serivce: GetSessionByAccessToken")
		}
	}

	return sessionByAccessToken, nil
}
