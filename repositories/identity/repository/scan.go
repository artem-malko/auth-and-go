package repository

import (
	"github.com/artem-malko/auth-and-go/infrastructure/database"
	"github.com/artem-malko/auth-and-go/models"
)

var fullIdentityFields = []string{
	idFieldName,
	accountIDFieldName,
	identityTypeFieldName,
	identityStatusFieldName,
	googleSocialIDFieldName,
	facebookSocialIDFieldName,
	emailFieldName,
	passwordHashFieldName,
	createdAtFieldName,
	updatedAtFieldName,
}

/**
Scan full identity info
*/
func (r *repository) scanFullIdentity(row database.RowScanner) (*models.Identity, error) {
	identity := new(dbIdentity)

	err := row.Scan(
		&identity.ID,
		&identity.AccountID,
		&identity.IdentityType,
		&identity.IdentityStatus,
		&identity.GoogleSocialID,
		&identity.FacebookSocialID,
		&identity.Email,
		&identity.PasswordHash,
		&identity.CreatedAt,
		&identity.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	baseIdentity := identity.GetBaseIdentity()

	return &baseIdentity, nil
}
