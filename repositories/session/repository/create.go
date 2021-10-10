package repository

import (
	"github.com/artem-malko/auth-and-go/infrastructure/database"

	"github.com/artem-malko/auth-and-go/models"
)

func (r *repository) CreateSession(executor database.QueryExecutor, session models.Session) error {
	queryBuilder := r.psql.Insert(sessionsTableName).Columns(
		sessionIDFieldName,
		accountIDFieldName,
		identityIDFieldName,
		clientIDFieldName,
		accessTokenFieldName,
		accessTokenExpiresDateFieldName,
		refreshTokenFieldName,
		refreshTokenExpiresDateFieldName,
	).Values(
		session.ID,
		session.AccountID,
		session.IdentityID,
		session.ClientID,
		session.AccessToken,
		session.AccessTokenExpiresDate,
		session.RefreshToken,
		session.RefreshTokenExpiresDate,
	)

	return r.exec(executor, queryBuilder)
}
