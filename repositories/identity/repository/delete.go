package repository

import (
	"github.com/artem-malko/auth-and-go/infrastructure/database"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

func (r *repository) DeleteIdentitiesByAccountID(
	executor database.QueryExecutor,
	accountID uuid.UUID,
) (err error) {
	queryBuilder := r.psql.Delete(identitiesTableName).
		Where(sq.Eq{accountIDFieldName: accountID})

	return r.exec(executor, queryBuilder)
}

func (r *repository) DeleteIdentitiesByIdentityIDs(
	executor database.QueryExecutor,
	identityIDs []uuid.UUID,
) error {
	queryBuilder := r.psql.Delete(identitiesTableName).
		Where(sq.Eq{idFieldName: identityIDs})

	return r.exec(executor, queryBuilder)
}
