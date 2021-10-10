package repository

import (
	"github.com/artem-malko/auth-and-go/infrastructure/database"
	"github.com/google/uuid"

	"github.com/artem-malko/auth-and-go/models"
)

func (r *repository) Create(
	executor database.QueryExecutor,
	token models.Token,
) (tokenID uuid.UUID, err error) {
	queryBuilder := r.psql.Insert(tokensTableName).Columns(
		tokenIDFieldName,
		tokenTypeFieldName,
		tokenStatusFieldName,
		accountIDFieldName,
		identityIDFieldName,
		clientIDFieldName,
		tokenExpiresDateFieldName,
	).Values(
		token.ID,
		token.TokenType,
		token.TokenStatus,
		token.AccountID,
		token.IdentityID,
		token.ClientID,
		token.ExpiresDate,
	)

	err = r.exec(executor, queryBuilder)

	if err != nil {
		return uuid.Nil, err
	}

	return token.ID, nil
}
