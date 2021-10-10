package repository

import (
	"database/sql"

	"github.com/artem-malko/auth-and-go/infrastructure/caller"

	"github.com/artem-malko/auth-and-go/infrastructure/database"

	"github.com/artem-malko/auth-and-go/repositories/token"

	sq "github.com/Masterminds/squirrel"
)

const tokensTableName = "tokens"
const repositoryPrefix = tokensTableName + " repository: "
const tokenIDFieldName = "id"
const tokenTypeFieldName = "token_type"
const tokenStatusFieldName = "token_status"
const accountIDFieldName = "account_id"
const identityIDFieldName = "identity_id"
const clientIDFieldName = "client_id"
const tokenExpiresDateFieldName = "expires_date"

type repository struct {
	psql sq.StatementBuilderType
}

// New creates postgresRepository
func New() (token.Repository, error) {
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

	return repositoryPrefix + "unknown caller"
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

func (r *repository) query(
	executor database.QueryExecutor, queryBuilder database.QueryBuilder, scanFunc func(rows *sql.Rows) error,
) error {
	return database.QueryWithCaller(
		getCaller(),
		executor,
		queryBuilder,
		scanFunc,
	)
}
