package repository

import (
	"database/sql"
	"strings"
	"time"

	"github.com/artem-malko/auth-and-go/infrastructure/database"

	sq "github.com/Masterminds/squirrel"
	"github.com/artem-malko/auth-and-go/models"
	"github.com/pkg/errors"
)

func (r *repository) DeleteUsedTokens(executor database.QueryExecutor) error {
	queryBuilder := r.psql.Delete(tokensTableName).
		Where(sq.Eq{tokenStatusFieldName: models.TokenStatusUsed})

	return r.exec(executor, queryBuilder)
}

func (r *repository) DeleteExpiredTokens(
	executor database.QueryExecutor,
	tokenType models.TokenType,
) ([]*models.Token, error) {
	queryBuilder := r.psql.Delete(tokensTableName).
		Where(sq.And{
			sq.Eq{tokenTypeFieldName: tokenType},
			sq.Eq{tokenStatusFieldName: models.TokenStatusActive},
			sq.LtOrEq{tokenExpiresDateFieldName: time.Now().Format(time.RFC3339)},
		}).
		Suffix("RETURNING " + strings.Join(fullTokenFields, ","))

	tokens := make([]*models.Token, 0)

	err := r.query(executor, queryBuilder, func(rows *sql.Rows) error {
		token, scanErr := r.scanFullToken(rows)

		if scanErr != nil {
			return errors.Wrap(scanErr, repositoryPrefix+"can`t scan info for DeleteExpiredTokens")
		}

		tokens = append(tokens, token)
		return nil
	})

	if err != nil {
		return nil, err
	}

	if len(tokens) == 0 {
		return nil, database.ErrRepositoryNoRowsAffected
	}

	return tokens, nil
}
