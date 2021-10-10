package repository

import (
	"database/sql"
	"strings"

	"github.com/artem-malko/auth-and-go/infrastructure/database"

	sq "github.com/Masterminds/squirrel"
	"github.com/artem-malko/auth-and-go/models"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

func (r *repository) UpdatePasswordHash(
	executor database.QueryExecutor,
	accountID uuid.UUID,
	passwordHash string,
) error {
	queryBuilder := r.psql.Update(identitiesTableName).
		Set(passwordHashFieldName, passwordHash).
		Where(sq.Eq{accountIDFieldName: accountID})

	return r.exec(executor, queryBuilder)
}

func (r *repository) UpdateIdentityStatus(
	executor database.QueryExecutor,
	identityID uuid.UUID,
	identityStatus models.IdentityStatus,
) (*models.Identity, error) {
	queryBuilder := r.psql.Update(identitiesTableName).
		Set(identityStatusFieldName, identityStatus).
		Where(sq.Eq{idFieldName: identityID}).
		Suffix("RETURNING " + strings.Join(fullIdentityFields, ","))

	row, err := r.queryRow(executor, queryBuilder)

	if err != nil {
		return nil, err
	}

	identity, err := r.scanFullIdentity(row)

	if err != nil {
		switch errors.Cause(err) {
		case sql.ErrNoRows:
			return nil, database.ErrRepositoryNoRowsAffected
		default:
			return nil, errors.Wrap(err, repositoryPrefix+"UpdateIdentityStatus exec query error")
		}
	}

	return identity, nil
}
