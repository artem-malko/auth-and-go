package service

import (
	"github.com/artem-malko/auth-and-go/infrastructure/database"
	"github.com/artem-malko/auth-and-go/services/token"
	"github.com/pkg/errors"

	"github.com/artem-malko/auth-and-go/models"
)

func (s *tokenService) DeleteUsedTokens(executor database.QueryExecutor) error {
	err := s.tokenRepository.DeleteUsedTokens(executor)

	if err != nil {
		return errors.Wrap(err, "token service: DeleteUsedTokens")
	}

	return nil
}

func (s *tokenService) DeleteExpiredTokens(executor database.QueryExecutor, tokenType models.TokenType) ([]*models.Token, error) {
	tokens, err := s.tokenRepository.DeleteExpiredTokens(executor, tokenType)

	if err != nil {
		switch errors.Cause(err) {
		case database.ErrRepositoryNoRowsAffected:
			return nil, token.ErrNoTokensUpdated
		default:
			return nil, errors.Wrap(err, "token service: DeleteExpiredTokens")
		}
	}

	return tokens, nil
}
