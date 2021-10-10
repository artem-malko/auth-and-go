package repository

import (
	"database/sql"

	"github.com/artem-malko/auth-and-go/infrastructure/database"

	sq "github.com/Masterminds/squirrel"

	"github.com/artem-malko/auth-and-go/models"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

func (r *repository) GetAccountByID(
	executor database.QueryExecutor,
	id uuid.UUID,
) (*models.Account, error) {
	queryBuilder := r.psql.Select(fullAccountFields...).
		From(accountsTableName).
		Where(sq.Eq{idFieldName: id})

	row, err := r.queryRow(executor, queryBuilder)

	if err != nil {
		return nil, err
	}

	u, err := r.scanFullAccount(row)

	if err != nil {
		switch errors.Cause(err) {
		case sql.ErrNoRows:
			return nil, database.ErrRepositoryNoEntitiesFound
		default:
			return nil, errors.Wrap(err, repositoryPrefix+"GetAccountByID scan error")
		}
	}

	return u, nil
}

func (r *repository) GetAccountsByIDsList(
	executor database.QueryExecutor,
	accountIDs []uuid.UUID,
) ([]*models.Account, error) {
	queryBuilder := r.psql.Select(fullAccountFields...).
		From(accountsTableName).
		Where(sq.Eq{idFieldName: accountIDs})

	foundAccounts := make([]*models.Account, 0)

	err := r.query(executor, queryBuilder, func(rows *sql.Rows) error {
		a, scanErr := r.scanFullAccount(rows)

		if scanErr != nil {
			return errors.Wrap(scanErr, repositoryPrefix+"can`t scan info for GetAccountsByIDsList")
		}

		if a != nil {
			foundAccounts = append(foundAccounts, a)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	if len(foundAccounts) == 0 || len(foundAccounts) != len(accountIDs) {
		return nil, database.ErrRepositoryNoEntitiesFound
	}

	return foundAccounts, nil
}

func (r *repository) GetAccountByName(
	executor database.QueryExecutor,
	accountName string,
) (*models.Account, error) {
	queryBuilder := r.psql.Select(fullAccountFields...).
		From(accountsTableName).
		Where(sq.Eq{accountNameFieldName: accountName})

	row, err := r.queryRow(executor, queryBuilder)

	if err != nil {
		return nil, err
	}

	u, err := r.scanFullAccount(row)

	if err != nil {
		switch errors.Cause(err) {
		case sql.ErrNoRows:
			return nil, database.ErrRepositoryNoEntitiesFound
		default:
			return nil, errors.Wrap(err, repositoryPrefix+"GetAccountByName scan error")
		}
	}

	return u, nil
}
