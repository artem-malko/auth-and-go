package repository

import (
	"database/sql"
	"strings"
	"time"

	"github.com/artem-malko/auth-and-go/infrastructure/database"

	"github.com/pkg/errors"

	sq "github.com/Masterminds/squirrel"

	"github.com/artem-malko/auth-and-go/models"
	"github.com/google/uuid"
)

func (r *repository) LoginAccount(
	executor database.QueryExecutor,
	accountID uuid.UUID,
	clientIP string,
) (*models.Account, error) {
	loginTime := time.Now()

	queryBuilder := r.psql.Update(accountsTableName).
		Set(lastIPFieldName, clientIP).
		Set(lastLoginFieldName, loginTime).
		Set(updatedAtFieldName, loginTime).
		Where(sq.Eq{idFieldName: accountID}).
		Suffix("RETURNING " + strings.Join(fullAccountFields, ","))

	row, err := r.queryRow(executor, queryBuilder)

	if err != nil {
		return nil, err
	}

	account, err := r.scanFullAccount(row)

	if err != nil {
		switch errors.Cause(err) {
		case sql.ErrNoRows:
			return nil, database.ErrRepositoryNoRowsAffected
		default:
			return nil, errors.Wrap(err, repositoryPrefix+"LoginAccount exec query error")
		}
	}

	return account, nil
}
