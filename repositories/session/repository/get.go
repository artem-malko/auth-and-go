package repository

import (
	"database/sql"

	"github.com/artem-malko/auth-and-go/infrastructure/database"

	sq "github.com/Masterminds/squirrel"
	"github.com/artem-malko/auth-and-go/models"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

func (r *repository) GetSessionByAccessToken(
	executor database.QueryExecutor,
	accessToken uuid.UUID,
) (*models.Session, error) {
	queryBuilder := r.psql.Select(fullSessionFields...).
		From(sessionsTableName).
		Where(sq.Eq{accessTokenFieldName: accessToken})

	row, err := r.queryRow(executor, queryBuilder)

	if err != nil {
		return nil, err
	}

	s, err := r.scanFullSession(row)

	if err != nil {
		switch errors.Cause(err) {
		case sql.ErrNoRows:
			return nil, database.ErrRepositoryNoEntitiesFound
		default:
			return nil, errors.Wrap(err, repositoryPrefix+"GetSessionByAccessToken scan error")
		}
	}

	return s, nil
}
