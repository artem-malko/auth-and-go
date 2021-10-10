package repository

import (
	"encoding/json"
	"time"

	"github.com/artem-malko/auth-and-go/infrastructure/database"

	"github.com/google/uuid"

	"github.com/artem-malko/auth-and-go/models"
	"github.com/pkg/errors"
)

func (r *repository) CreateAccount(
	executor database.QueryExecutor,
	account models.Account,
) (*models.Account, error) {
	dbAccount := createDBAccount(account)

	profileData, err := json.Marshal(dbAccount.Profile)

	if err != nil {
		return nil, errors.Wrap(err, repositoryPrefix+"CreateAccount marshal profile error")
	}

	settingsData, err := json.Marshal(dbAccount.Settings)

	if err != nil {
		return nil, errors.Wrap(err, repositoryPrefix+"CreateAccount marshal settings error")
	}

	queryBuilder := r.psql.Insert(accountsTableName).Columns(
		idFieldName,
		accountTypeFieldName,
		accountStatusFieldName,
		accountNameFieldName,
		profileFieldName,
		settingsFieldName,
		lastIPFieldName,
	).Values(
		dbAccount.ID,
		dbAccount.AccountType,
		dbAccount.AccountStatus,
		dbAccount.AccountName,
		profileData,
		settingsData,
		dbAccount.LastIP,
	).Suffix(
		"RETURNING " +
			idFieldName + ", " +
			lastLoginFieldName + ", " +
			createdAtFieldName + ", " +
			updatedAtFieldName,
	)

	row, err := r.queryRow(executor, queryBuilder)

	if err != nil {
		return nil, err
	}

	var id uuid.UUID
	var lastLogin, createdAt, updatedAt time.Time

	err = row.Scan(&id, &lastLogin, &createdAt, &updatedAt)

	// @TODO handle pkey error
	if err != nil {
		return nil, errors.Wrap(err, repositoryPrefix+"CreateAccount QueryRow error")
	}

	var createdAccount = dbAccount.GetBaseAccount()
	// Won't be zero time, cause DB will fill these dates
	createdAccount.CreatedAt = createdAt
	createdAccount.UpdatedAt = updatedAt
	createdAccount.LastLogin = lastLogin
	createdAccount.ID = id

	return &createdAccount, nil
}
