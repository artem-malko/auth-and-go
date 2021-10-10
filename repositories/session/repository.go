package session

import (
	"time"

	"github.com/artem-malko/auth-and-go/infrastructure/database"

	"github.com/artem-malko/auth-and-go/models"
	"github.com/google/uuid"
)

// Repository is an interface for Session repository
type Repository interface {
	CreateSession(executor database.QueryExecutor, session models.Session) error
	UpdateSessionByRefreshToken(
		executor database.QueryExecutor,
		refreshToken uuid.UUID,
		accessTokenExpiresDate, refreshTokenExpiresDate time.Time,
	) (*models.Session, error)
	DeleteAllSessionsByAccountID(executor database.QueryExecutor, accountID uuid.UUID) error
	DeleteSessionBySessionID(executor database.QueryExecutor, sessionID uuid.UUID) error
	GetSessionByAccessToken(executor database.QueryExecutor, accessToken uuid.UUID) (*models.Session, error)
	DeleteExpiredSessions(executor database.QueryExecutor) error
}
