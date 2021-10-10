package repository

import (
	"encoding/json"

	"github.com/artem-malko/auth-and-go/infrastructure/database"

	"github.com/pkg/errors"

	"github.com/artem-malko/auth-and-go/models"
)

var fullAccountFields = []string{
	idFieldName,
	accountTypeFieldName,
	accountStatusFieldName,
	accountNameFieldName,
	profileFieldName,
	settingsFieldName,
	lastIPFieldName,
	lastLoginFieldName,
	createdAtFieldName,
	updatedAtFieldName,
}

/**
Scan full user info
*/
func (r *repository) scanFullAccount(row database.RowScanner) (*models.Account, error) {
	dbAccount := new(dbAccount)
	parsedProfile := new(dbProfile)
	parsedSettings := new(dbSettings)

	var profile []byte
	var settings []byte

	err := row.Scan(
		&dbAccount.ID,
		&dbAccount.AccountType,
		&dbAccount.AccountStatus,
		&dbAccount.AccountName,
		&profile,
		&settings,
		&dbAccount.LastIP,
		&dbAccount.LastLogin,
		&dbAccount.CreatedAt,
		&dbAccount.UpdatedAt,
	)

	if err != nil {
		return nil, errors.Wrap(err, "account repository: full scan error")
	}

	err = json.Unmarshal(profile, &parsedProfile)

	if err != nil {
		return nil, errors.Wrap(err, "account repository: unmarshal profile error")
	}

	err = json.Unmarshal(settings, &parsedSettings)

	if err != nil {
		return nil, errors.Wrap(err, "account repository: unmarshal settings error")
	}

	dbAccount.Profile = *parsedProfile
	dbAccount.Settings = *parsedSettings

	baseAccount := dbAccount.GetBaseAccount()

	return &baseAccount, nil
}
