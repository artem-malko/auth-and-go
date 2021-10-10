package repository

import (
	"github.com/artem-malko/auth-and-go/infrastructure/database"
	"github.com/artem-malko/auth-and-go/models"
)

var fullSessionFields = []string{
	sessionIDFieldName,
	accountIDFieldName,
	identityIDFieldName,
	clientIDFieldName,
	accessTokenFieldName,
	accessTokenExpiresDateFieldName,
	refreshTokenFieldName,
	refreshTokenExpiresDateFieldName,
}

/**
Scan full session info
*/
func (r *repository) scanFullSession(row database.RowScanner) (*models.Session, error) {
	session := new(models.Session)

	err := row.Scan(
		&session.ID,
		&session.AccountID,
		&session.IdentityID,
		&session.ClientID,
		&session.AccessToken,
		&session.AccessTokenExpiresDate,
		&session.RefreshToken,
		&session.RefreshTokenExpiresDate,
	)

	if err != nil {
		return nil, err
	}

	return session, nil
}
