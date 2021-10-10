package repository

import (
	"database/sql"
	"strings"
	"time"

	"github.com/artem-malko/auth-and-go/infrastructure/database"
	"github.com/artem-malko/auth-and-go/repositories/account"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"

	sq "github.com/Masterminds/squirrel"
	"github.com/artem-malko/auth-and-go/models"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

func (r *repository) ConfirmAccount(
	executor database.QueryExecutor,
	accountID uuid.UUID,
) (*models.Account, error) {
	// @TODO add settings update
	setMap := map[string]interface{}{
		accountStatusFieldName: models.AccountStatusConfirmed,
		updatedAtFieldName:     time.Now(),
	}

	queryBuilder := r.psql.Update(accountsTableName).
		SetMap(setMap).
		Where(sq.Eq{idFieldName: accountID}).
		Suffix("RETURNING " + strings.Join(fullAccountFields, ","))

	row, err := r.queryRow(executor, queryBuilder)

	if err != nil {
		return nil, err
	}

	a, err := r.scanFullAccount(row)

	if err != nil {
		switch errors.Cause(err) {
		case sql.ErrNoRows:
			return nil, database.ErrRepositoryNoRowsAffected
		default:
			return nil, errors.Wrap(err, repositoryPrefix+"ConfirmAccount exec query error")
		}
	}

	return a, nil
}

func (r *repository) UpdateAccountName(
	executor database.QueryExecutor,
	accountID uuid.UUID,
	name string,
) error {
	setMap := map[string]interface{}{
		accountNameFieldName: name,
		updatedAtFieldName:   time.Now(),
	}

	queryBuilder := r.psql.Update(accountsTableName).
		SetMap(setMap).
		Where(sq.Eq{idFieldName: accountID})

	err := r.exec(executor, queryBuilder)

	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.UniqueViolation && pgErr.ConstraintName == accountNameUniqIndexName {
				return account.ErrRepositoryAccountNameConstraint
			}

			return errors.Wrap(err, repositoryPrefix+"PGX error")
		}

		return err
	}

	return nil
}
