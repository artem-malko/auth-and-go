package repository

import (
	"time"

	"github.com/artem-malko/auth-and-go/infrastructure/database"

	"github.com/google/uuid"

	sq "github.com/Masterminds/squirrel"
	"github.com/artem-malko/auth-and-go/models"
)

// DeactivateAccountByID set status "Deleted" for account and drop all account data
func (r *repository) DeactivateAccountByID(executor database.QueryExecutor, accountID uuid.UUID) error {
	valuesToUpdate := map[string]interface{}{}

	valuesToUpdate[accountStatusFieldName] = models.AccountStatusDeleted
	valuesToUpdate[accountNameFieldName] = "Anonymous"
	// @TODO remove settings
	valuesToUpdate[updatedAtFieldName] = time.Now()

	queryBuilder := r.psql.Update(accountsTableName).
		SetMap(valuesToUpdate).
		Where(sq.And{sq.Eq{idFieldName: accountID}, sq.NotEq{accountStatusFieldName: models.AccountStatusDeleted}}).
		Suffix("")

	return r.exec(executor, queryBuilder)
}

func (r *repository) DeleteUnconfirmedAccountsByAccountIDs(
	executor database.QueryExecutor,
	accountIDs []uuid.UUID,
) error {
	queryBuilder := r.psql.Delete(accountsTableName).
		Where(sq.And{
			sq.Eq{idFieldName: accountIDs},
			sq.Eq{accountStatusFieldName: models.AccountStatusUnconfirmed},
		})

	return r.exec(executor, queryBuilder)
}
