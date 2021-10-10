package token

import (
	"github.com/artem-malko/auth-and-go/infrastructure/database"

	"github.com/artem-malko/auth-and-go/models"
	"github.com/google/uuid"
)

// IdentityRepository is an interface for any Repository
type Repository interface {
	Create(executor database.QueryExecutor, token models.Token) (tokenID uuid.UUID, err error)
	GetTokensByIdentityID(
		executor database.QueryExecutor,
		identityID uuid.UUID,
		tokenType models.TokenType,
	) ([]*models.Token, error)
	UpdateStatus(
		executor database.QueryExecutor,
		tokenID uuid.UUID,
		tokenStatus models.TokenStatus,
	) (*models.Token, error)
	DeleteUsedTokens(executor database.QueryExecutor) error
	DeleteExpiredTokens(
		executor database.QueryExecutor,
		tokenType models.TokenType,
	) ([]*models.Token, error)
}
