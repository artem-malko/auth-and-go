package service

import (
	"time"

	"github.com/pkg/errors"

	"github.com/artem-malko/auth-and-go/constants"
	"github.com/artem-malko/auth-and-go/infrastructure/database"

	"github.com/artem-malko/auth-and-go/models"
	"github.com/google/uuid"
)

func (s *tokenService) Create(
	executor database.QueryExecutor,
	tokenType models.TokenType,
	clientID models.ClientID,
	accountID, identityID uuid.UUID,
) (*models.Token, error) {
	var expiresDate time.Time

	switch tokenType {
	case models.TokenTypeRegistrationConfirmation:
		expiresDate = time.Now().
			Add(time.Second * time.Duration(constants.Values.RegistrationConfirmationTokenMaxAgeInSeconds))
	case models.TokenTypeAutoLogin:
		expiresDate = time.Now().
			Add(time.Second * time.Duration(constants.Values.AutoLoginTokenMaxAgeInSeconds))
	case models.TokenTypeEmailConfirmation:
		expiresDate = time.Now().
			Add(time.Second * time.Duration(constants.Values.ChangeEmailConfirmationTokenMaxAgeInSeconds))
	}

	token := models.Token{
		ID:          uuid.New(),
		TokenType:   tokenType,
		TokenStatus: models.TokenStatusActive,
		AccountID:   accountID,
		IdentityID:  identityID,
		ClientID:    clientID,
		ExpiresDate: expiresDate,
	}

	_, err := s.tokenRepository.Create(executor, token)

	if err != nil {
		return nil, errors.Wrap(err, "token service: Create error")
	}

	return &token, nil
}
