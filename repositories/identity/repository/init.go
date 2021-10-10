package repository

import (
	"database/sql"
	"time"

	"github.com/artem-malko/auth-and-go/infrastructure/caller"

	"github.com/artem-malko/auth-and-go/infrastructure/database"

	"gopkg.in/guregu/null.v3"

	"github.com/artem-malko/auth-and-go/infrastructure/convert"
	"github.com/artem-malko/auth-and-go/models"
	"github.com/artem-malko/auth-and-go/repositories/identity"
	"github.com/google/uuid"

	sq "github.com/Masterminds/squirrel"
)

const identitiesTableName = "identities"
const repositoryPrefix = identitiesTableName + " repository: "
const idFieldName = "id"
const accountIDFieldName = "account_id"
const identityTypeFieldName = "identity_type"
const identityStatusFieldName = "identity_status"
const googleSocialIDFieldName = "google_social_id"
const facebookSocialIDFieldName = "facebook_social_id"
const emailFieldName = "email"
const passwordHashFieldName = "password_hash"
const createdAtFieldName = "created_at"
const updatedAtFieldName = "updated_at"

const emailUniqIndexName = "identities_email_uniq_idx"
const googleSocialIDUniqIndexName = "identities_google_social_id_uniq_idx"
const facebookSocialIDUniqIndexName = "identities_facebook_social_id_uniq_idx"

type dbIdentity struct {
	ID               uuid.UUID
	AccountID        uuid.UUID
	IdentityType     models.IdentityType
	IdentityStatus   models.IdentityStatus
	GoogleSocialID   null.String
	FacebookSocialID null.String
	Email            null.String
	PasswordHash     null.String
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

func createDBIdentity(i models.Identity) dbIdentity {
	return dbIdentity{
		ID:               i.ID,
		AccountID:        i.AccountID,
		IdentityType:     i.IdentityType,
		IdentityStatus:   i.IdentityStatus,
		GoogleSocialID:   convert.NewSQLNullString(i.GoogleSocialID),
		FacebookSocialID: convert.NewSQLNullString(i.FacebookSocialID),
		Email:            convert.NewSQLNullString(i.Email),
		PasswordHash:     convert.NewSQLNullString(i.PasswordHash),
		CreatedAt:        i.CreatedAt,
		UpdatedAt:        i.UpdatedAt,
	}
}

func (d *dbIdentity) GetBaseIdentity() models.Identity {
	return models.Identity{
		ID:               d.ID,
		AccountID:        d.AccountID,
		IdentityType:     d.IdentityType,
		IdentityStatus:   d.IdentityStatus,
		GoogleSocialID:   d.GoogleSocialID.String,
		FacebookSocialID: d.FacebookSocialID.String,
		Email:            d.Email.String,
		PasswordHash:     d.PasswordHash.String,
		CreatedAt:        d.CreatedAt,
		UpdatedAt:        d.UpdatedAt,
	}
}

// New creates postgresRepository
func New() (identity.Repository, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	return &repository{
		psql: psql,
	}, nil
}

type repository struct {
	psql sq.StatementBuilderType
}

func getCaller() string {
	// Remove getCaller from stack
	if callerName, ok := caller.GetCaller(2); ok == true {
		return callerName
	}

	return "identity repository: unknown caller"
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
