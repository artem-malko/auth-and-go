package token

import (
	"errors"

	"github.com/artem-malko/auth-and-go/infrastructure/database"
	"github.com/artem-malko/auth-and-go/models"
	"github.com/google/uuid"
)

var (
	ErrNoTokensUpdated = errors.New("tokens with sent params are not updated")
	ErrNoTokens        = errors.New("tokens with sent params are not found")
)

type Repository interface {
	Create(executor database.QueryExecutor, token models.Token) (tokenID uuid.UUID, err error)
	GetTokensByIdentityID(
		executor database.QueryExecutor,
		identityID uuid.UUID,
		tokenType models.TokenType,
	) ([]*models.Token, error)
	UpdateStatus(
		executor database.QueryExecutor,
		TokenID uuid.UUID,
		TokenStatus models.TokenStatus,
	) (*models.Token, error)
	DeleteUsedTokens(executor database.QueryExecutor) error
	DeleteExpiredTokens(executor database.QueryExecutor, tokenType models.TokenType) ([]*models.Token, error)
}

type Service interface {
	Create(
		executor database.QueryExecutor,
		tokenType models.TokenType,
		clientID models.ClientID,
		accountID, identityID uuid.UUID,
	) (*models.Token, error)
	GetActiveTokenByIdentityID(executor database.QueryExecutor, identityID uuid.UUID, tokenType models.TokenType) (*models.Token, error)
	Use(executor database.QueryExecutor, tokenID uuid.UUID) (*models.Token, error)
	DeleteUsedTokens(executor database.QueryExecutor) error
	DeleteExpiredTokens(executor database.QueryExecutor, tokenType models.TokenType) ([]*models.Token, error)
}
