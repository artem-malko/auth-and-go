package service

import (
	"time"

	"github.com/artem-malko/auth-and-go/infrastructure/database"
	"github.com/artem-malko/auth-and-go/models"
	"github.com/artem-malko/auth-and-go/services/token"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

func (s *tokenService) GetActiveTokenByIdentityID(executor database.QueryExecutor, identityID uuid.UUID, tokenType models.TokenType) (*models.Token, error) {
	tokens, err := s.tokenRepository.GetTokensByIdentityID(executor, identityID, tokenType)

	if err != nil {
		switch errors.Cause(err) {
		case database.ErrRepositoryNoEntitiesFound:
			return nil, token.ErrNoTokens
		default:
			return nil, errors.Wrap(err, "token service: GetActiveTokenByIdentityID")
		}
	}

	var activeToken *models.Token

	for _, t := range tokens {
		if t.TokenStatus == models.TokenStatusActive && time.Now().Before(t.ExpiresDate) {
			activeToken = t
			break
		}
	}

	if activeToken == nil {
		return nil, token.ErrNoTokens
	}

	return activeToken, nil
}
