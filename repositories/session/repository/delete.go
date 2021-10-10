package repository

import (
	"time"

	"github.com/artem-malko/auth-and-go/infrastructure/database"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

func (r *repository) DeleteAllSessionsByAccountID(
	executor database.QueryExecutor,
	accountID uuid.UUID,
) error {
	queryBuilder := r.psql.Delete(sessionsTableName).
		Where(sq.Eq{accountIDFieldName: accountID})

	return r.exec(executor, queryBuilder)
}

func (r *repository) DeleteSessionBySessionID(
	executor database.QueryExecutor,
	sessionID uuid.UUID,
) error {
	queryBuilder := r.psql.Delete(sessionsTableName).
		Where(sq.Eq{sessionIDFieldName: sessionID})

	return r.exec(executor, queryBuilder)
}

func (r *repository) DeleteExpiredSessions(executor database.QueryExecutor) error {
	queryBuilder := r.psql.Delete(sessionsTableName).
		Where(sq.Lt{refreshTokenExpiresDateFieldName: time.Now().Format(time.RFC3339)})

	return r.exec(executor, queryBuilder)
}
