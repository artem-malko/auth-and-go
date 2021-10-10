package repository

import (
	"database/sql"

	"github.com/artem-malko/auth-and-go/infrastructure/caller"

	"github.com/artem-malko/auth-and-go/infrastructure/database"

	sq "github.com/Masterminds/squirrel"

	"github.com/artem-malko/auth-and-go/repositories/session"
)

const sessionsTableName = "sessions"
const repositoryPrefix = sessionsTableName + " repository: "
const sessionIDFieldName = "id"
const accountIDFieldName = "account_id"
const identityIDFieldName = "identity_id"
const clientIDFieldName = "client_id"
const accessTokenFieldName = "access_token"
const accessTokenExpiresDateFieldName = "access_token_expires_date"
const refreshTokenFieldName = "refresh_token"
const refreshTokenExpiresDateFieldName = "refresh_token_expires_date"

type repository struct {
	db   *sql.DB
	psql sq.StatementBuilderType
}

// New creates postgresRepository
func New() (session.Repository, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	return &repository{
		psql: psql,
	}, nil
}

func getCaller() string {
	// Remove getCaller from stack
	if callerName, ok := caller.GetCaller(2); ok == true {
		return callerName
	}

	return "session repository: unknown caller"
}

func (r *repository) exec(executor database.QueryExecutor, queryBuilder database.QueryBuilder) error {
	return database.ExecWithCaller(
		getCaller(),
		executor,
		queryBuilder,
	)
}

func (r *repository) queryRow(executor database.QueryExecutor, queryBuilder database.QueryBuilder) (*sql.Row, error) {
	return database.QueryRowWithCaller(
		getCaller(),
		executor,
		queryBuilder,
	)
}
