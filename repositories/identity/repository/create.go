package repository

import (
	"time"

	"github.com/artem-malko/auth-and-go/repositories/identity"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"

	"github.com/artem-malko/auth-and-go/infrastructure/database"

	"github.com/artem-malko/auth-and-go/models"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

func (r *repository) CreateEmailIdentity(
	executor database.QueryExecutor,
	identityToCreate models.Identity,
) (*models.Identity, error) {
	DBIdentity := createDBIdentity(identityToCreate)

	queryBuilder := r.psql.Insert(identitiesTableName).Columns(
		idFieldName,
		accountIDFieldName,
		identityTypeFieldName,
		identityStatusFieldName,
		emailFieldName,
		passwordHashFieldName,
	).Values(
		DBIdentity.ID,
		DBIdentity.AccountID,
		DBIdentity.IdentityType,
		DBIdentity.IdentityStatus,
		DBIdentity.Email,
		DBIdentity.PasswordHash,
	).Suffix("RETURNING " + idFieldName + ", " + createdAtFieldName + ", " + updatedAtFieldName)

	row, err := r.queryRow(executor, queryBuilder)

	if err != nil {
		return nil, err
	}

	var id uuid.UUID
	var createdAt, updatedAt time.Time

	err = row.Scan(&id, &createdAt, &updatedAt)

	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.UniqueViolation && pgErr.ConstraintName == emailUniqIndexName {
				return nil, identity.ErrRepositoryEmailConstraint
			}

			return nil, errors.Wrap(pgErr, repositoryPrefix+"CreateEmailIdentity exec query PGX error")
		}

		return nil, errors.Wrap(err, repositoryPrefix+"CreateEmailIdentity exec query error")
	}

	var createdIdentity = DBIdentity.GetBaseIdentity()
	// Won't be zero time, cause DB will fill these dates
	createdIdentity.CreatedAt = createdAt
	createdIdentity.UpdatedAt = updatedAt
	createdIdentity.ID = id

	return &createdIdentity, nil
}

func (r *repository) CreateSocialIdentity(
	executor database.QueryExecutor,
	identityToCreate models.Identity,
) (*models.Identity, error) {
	DBIdentity := createDBIdentity(identityToCreate)

	socialIDFieldNameAndValue, err := getSocialIDFieldNameAndValue(identityToCreate)

	if err != nil {
		return nil, err
	}

	queryBuilder := r.psql.Insert(identitiesTableName).Columns(
		idFieldName,
		accountIDFieldName,
		identityTypeFieldName,
		identityStatusFieldName,
		emailFieldName,
		socialIDFieldNameAndValue.FieldName,
	).Values(
		DBIdentity.ID,
		DBIdentity.AccountID,
		DBIdentity.IdentityType,
		DBIdentity.IdentityStatus,
		DBIdentity.Email,
		socialIDFieldNameAndValue.Value,
	).Suffix("RETURNING " + idFieldName + ", " + createdAtFieldName + ", " + updatedAtFieldName)

	row, err := r.queryRow(executor, queryBuilder)

	if err != nil {
		return nil, err
	}

	var id uuid.UUID
	var createdAt, updatedAt time.Time

	err = row.Scan(&id, &createdAt, &updatedAt)

	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.UniqueViolation {
				switch pgErr.ConstraintName {
				case googleSocialIDUniqIndexName:
					return nil, identity.ErrRepositorySocialGoogleConstraint
				case facebookSocialIDUniqIndexName:
					return nil, identity.ErrRepositorySocialFacebookConstraint
				}
			}

			return nil, errors.Wrap(pgErr, repositoryPrefix+"CreateSocialIdentity exec query PGX error")
		}

		return nil, errors.Wrap(err, repositoryPrefix+"CreateSocialIdentity exec query error")
	}

	var createdIdentity = DBIdentity.GetBaseIdentity()
	// Won't be zero time, cause DB will fill these dates
	createdIdentity.CreatedAt = createdAt
	createdIdentity.UpdatedAt = updatedAt
	createdIdentity.ID = id

	return &createdIdentity, nil
}

type socialIDFieldNameAndValue struct {
	FieldName string
	Value     string
}

func getSocialIDFieldNameAndValue(
	identityToCreate models.Identity,
) (*socialIDFieldNameAndValue, error) {
	switch identityToCreate.IdentityType {
	case "google":
		return &socialIDFieldNameAndValue{
			FieldName: googleSocialIDFieldName,
			Value:     identityToCreate.GoogleSocialID,
		}, nil
	case "facebook":
		return &socialIDFieldNameAndValue{
			FieldName: facebookSocialIDFieldName,
			Value:     identityToCreate.FacebookSocialID,
		}, nil
	default:
		return nil, identity.ErrRepositoryUnknownSocialNetworkType
	}
}
