package repository

import (
	"database/sql"
	"strings"
	"time"

	"github.com/artem-malko/auth-and-go/infrastructure/database"

	sq "github.com/Masterminds/squirrel"
	"github.com/pkg/errors"

	"github.com/artem-malko/auth-and-go/models"
	"github.com/google/uuid"
)

func (r *repository) UpdateStatus(
	executor database.QueryExecutor,
	tokenID uuid.UUID,
	tokenStatus models.TokenStatus,
) (*models.Token, error) {
	queryBuilder := r.psql.Update(tokensTableName).
		Set(tokenStatusFieldName, tokenStatus).
		Where(sq.And{
			sq.Eq{tokenIDFieldName: tokenID},
			sq.Eq{tokenStatusFieldName: models.TokenStatusActive},
			sq.Gt{tokenExpiresDateFieldName: time.Now().Format(time.RFC3339)},
		}).
		Suffix("RETURNING " + strings.Join(fullTokenFields, ","))

	row, err := r.queryRow(executor, queryBuilder)

	if err != nil {
		return nil, err
	}

	identity, err := r.scanFullToken(row)

	if err != nil {
		switch errors.Cause(err) {
		case sql.ErrNoRows:
			return nil, database.ErrRepositoryNoRowsAffected
		default:
			return nil, errors.Wrap(err, repositoryPrefix+"UpdateStatus exec query error")
		}
	}

	return identity, nil
}
