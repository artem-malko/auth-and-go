package repository

import (
	"github.com/artem-malko/auth-and-go/infrastructure/database"
	"github.com/artem-malko/auth-and-go/models"
	"github.com/pkg/errors"
)

var fullTokenFields = []string{
	tokenIDFieldName,
	tokenTypeFieldName,
	tokenStatusFieldName,
	accountIDFieldName,
	identityIDFieldName,
	clientIDFieldName,
	tokenExpiresDateFieldName,
}

/**
Scan full token info
*/
func (r *repository) scanFullToken(row database.RowScanner) (*models.Token, error) {
	token := new(models.Token)

	err := row.Scan(
		&token.ID,
		&token.TokenType,
		&token.TokenStatus,
		&token.AccountID,
		&token.IdentityID,
		&token.ClientID,
		&token.ExpiresDate,
	)

	if err != nil {
		return nil, errors.Wrap(err, repositoryPrefix+"scanFullToken err")
	}

	return token, nil
}
