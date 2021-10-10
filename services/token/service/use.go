package service

import (
	"github.com/artem-malko/auth-and-go/infrastructure/database"
	"github.com/artem-malko/auth-and-go/services/token"

	"github.com/pkg/errors"

	"github.com/artem-malko/auth-and-go/models"
	"github.com/google/uuid"
)

func (s *tokenService) Use(executor database.QueryExecutor, tokenID uuid.UUID) (*models.Token, error) {
	usedToken, err := s.tokenRepository.UpdateStatus(executor, tokenID, models.TokenStatusUsed)

	if err != nil {
		switch errors.Cause(err) {
		case database.ErrRepositoryNoRowsAffected:
			return nil, token.ErrNoTokensUpdated
		default:
			return nil, errors.Wrap(err, "token service: Use token err")
		}
	}

	return usedToken, nil
}
