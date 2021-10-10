package repository

import (
	"database/sql"
	"strings"
	"time"

	"github.com/artem-malko/auth-and-go/infrastructure/database"

	sq "github.com/Masterminds/squirrel"
	"github.com/artem-malko/auth-and-go/models"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

func (r *repository) UpdateSessionByRefreshToken(
	executor database.QueryExecutor,
	refreshToken uuid.UUID,
	accessTokenExpiresDate, refreshTokenExpiresDate time.Time,
) (*models.Session, error) {
	queryBuilder := r.psql.Update(sessionsTableName).
		Set(accessTokenExpiresDateFieldName, accessTokenExpiresDate).
		Set(refreshTokenExpiresDateFieldName, refreshTokenExpiresDate).
		Where(sq.And{
			sq.Eq{refreshTokenFieldName: refreshToken},
			sq.Gt{refreshTokenExpiresDateFieldName: time.Now()},
		}).
		Suffix("RETURNING " + strings.Join(fullSessionFields, ","))

	row, err := r.queryRow(executor, queryBuilder)

	if err != nil {
		return nil, err
	}

	s, err := r.scanFullSession(row)

	if err != nil {
		switch errors.Cause(err) {
		case sql.ErrNoRows:
			return nil, database.ErrRepositoryNoRowsAffected
		default:
			return nil, errors.Wrap(err, repositoryPrefix+"UpdateSessionByRefreshToken exec query error")
		}
	}

	return s, nil
}
