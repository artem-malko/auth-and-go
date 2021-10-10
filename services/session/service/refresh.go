package service

import (
	"time"

	"github.com/artem-malko/auth-and-go/infrastructure/database"
	"github.com/artem-malko/auth-and-go/services/session"
	"github.com/pkg/errors"

	"github.com/artem-malko/auth-and-go/constants"
	"github.com/artem-malko/auth-and-go/models"
	"github.com/google/uuid"
)

func (s *sessionService) RefreshSession(
	executor database.QueryExecutor,
	refreshToken uuid.UUID,
) (*models.Session, error) {
	accessTokenExpiresDate := time.Now().
		Add(time.Second * time.Duration(constants.Values.AccessTokenMaxAgeInSeconds))
	refreshTokenExpiresDate := time.Now().
		Add(time.Second * time.Duration(constants.Values.RefreshTokenMaxAgeInSeconds))

	refreshedSession, err := s.sessionRepository.UpdateSessionByRefreshToken(
		executor,
		refreshToken,
		accessTokenExpiresDate,
		refreshTokenExpiresDate,
	)

	if err != nil {
		switch errors.Cause(err) {
		case database.ErrRepositoryNoRowsAffected:
			return nil, session.ErrNoSessionsUpdated
		default:
			return nil, errors.Wrap(err, "session service: RefreshSession")
		}
	}

	return refreshedSession, nil
}
