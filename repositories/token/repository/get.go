package repository

import (
	"database/sql"

	"github.com/artem-malko/auth-and-go/infrastructure/database"

	sq "github.com/Masterminds/squirrel"
	"github.com/artem-malko/auth-and-go/models"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

func (r *repository) GetTokensByIdentityID(
	executor database.QueryExecutor,
	identityID uuid.UUID,
	tokenType models.TokenType,
) ([]*models.Token, error) {
	queryBuilder := r.psql.Select(fullTokenFields...).
		From(tokensTableName).
		Where(sq.And{
			sq.Eq{identityIDFieldName: identityID},
			sq.Eq{tokenTypeFieldName: tokenType},
		})

	tokens := make([]*models.Token, 0)

	err := r.query(executor, queryBuilder, func(rows *sql.Rows) error {
		token, scanErr := r.scanFullToken(rows)

		if scanErr != nil {
			return errors.Wrap(scanErr, repositoryPrefix+"can`t scan info for GetTokensByIdentityID")
		}

		tokens = append(tokens, token)
		return nil
	})

	if err != nil {
		return nil, err
	}

	if len(tokens) == 0 {
		return nil, database.ErrRepositoryNoEntitiesFound
	}

	return tokens, nil
}
