package session

import (
	"errors"
	"time"

	"github.com/artem-malko/auth-and-go/infrastructure/database"
	"github.com/artem-malko/auth-and-go/models"
	"github.com/google/uuid"
)

var (
	ErrSessionNotFound   = errors.New("session not found")
	ErrNoSessionsUpdated = errors.New("sessions with sent params are not updated")
)

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

type Service interface {
	GetSessionByAccessToken(executor database.QueryExecutor, accessToken uuid.UUID) (*models.Session, error)
	CreateSession(
		executor database.QueryExecutor,
		accountID, identityID uuid.UUID,
		clientID models.ClientID,
	) (*models.Session, error)
	DeleteAllSessionsByAccountID(executor database.QueryExecutor, accountID uuid.UUID) error
	DeleteSessionBySessionID(executor database.QueryExecutor, sessionID uuid.UUID) error
	DeleteExpiredSessions(executor database.QueryExecutor) error
	RefreshSession(executor database.QueryExecutor, refreshToken uuid.UUID) (*models.Session, error)
}
