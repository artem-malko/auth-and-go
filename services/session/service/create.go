package service

import (
	"time"

	"github.com/pkg/errors"

	"github.com/artem-malko/auth-and-go/constants"
	"github.com/artem-malko/auth-and-go/infrastructure/database"

	"github.com/artem-malko/auth-and-go/models"
	"github.com/google/uuid"
)

func (s *sessionService) CreateSession(
	executor database.QueryExecutor,
	accountID, identityID uuid.UUID,
	clientID models.ClientID,
) (*models.Session, error) {
	session := models.Session{
		AccountID:   accountID,
		IdentityID:  identityID,
		ID:          uuid.New(),
		ClientID:    clientID,
		AccessToken: uuid.New(),
		AccessTokenExpiresDate: time.Now().
			Add(time.Second * time.Duration(constants.Values.AccessTokenMaxAgeInSeconds)),
		RefreshToken: uuid.New(),
		RefreshTokenExpiresDate: time.Now().
			Add(time.Second * time.Duration(constants.Values.RefreshTokenMaxAgeInSeconds)),
	}

	err := s.sessionRepository.CreateSession(executor, session)

	if err != nil {
		return nil, errors.Wrap(err, "session service: CreateSession error")
	}

	return &session, nil
}
