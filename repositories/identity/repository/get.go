package repository

import (
	"database/sql"

	"github.com/artem-malko/auth-and-go/infrastructure/database"

	"github.com/google/uuid"

	sq "github.com/Masterminds/squirrel"
	"github.com/artem-malko/auth-and-go/models"
	"github.com/pkg/errors"
)

func (r *repository) GetIdentitiesByEmail(
	executor database.QueryExecutor,
	email string,
) ([]*models.Identity, error) {
	queryBuilder := r.psql.Select(fullIdentityFields...).
		From(identitiesTableName).
		Where(sq.And{sq.Eq{emailFieldName: email}})

	identities := make([]*models.Identity, 0)

	err := r.query(executor, queryBuilder, func(rows *sql.Rows) error {
		identity, scanErr := r.scanFullIdentity(rows)

		if scanErr != nil {
			return errors.Wrap(scanErr, repositoryPrefix+"can`t scan info for GetIdentitiesByEmail")
		}

		identities = append(identities, identity)

		return nil
	})

	if err != nil {
		return nil, err
	}

	if len(identities) == 0 {
		return nil, database.ErrRepositoryNoEntitiesFound
	}

	return identities, nil
}

func (r *repository) GetEmailIdentityByEmail(
	executor database.QueryExecutor,
	email string,
) (*models.Identity, error) {
	queryBuilder := r.psql.Select(fullIdentityFields...).
		From(identitiesTableName).
		Where(sq.And{sq.Eq{emailFieldName: email}, sq.Eq{identityTypeFieldName: "email"}})

	row, err := r.queryRow(executor, queryBuilder)

	if err != nil {
		return nil, err
	}

	i, err := r.scanFullIdentity(row)

	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, database.ErrRepositoryNoEntitiesFound
		}

		return nil, errors.Wrap(err, repositoryPrefix+"GetEmailIdentityByEmail scan error")
	}

	return i, nil
}

func (r *repository) GetIdentityBySocialID(
	executor database.QueryExecutor,
	socialID string,
	socialNetworkType models.SocialNetworkType,
) (*models.Identity, error) {
	var socialNetworkIDFieldName string

	if socialNetworkType == models.SocialNetworkTypeFacebook {
		socialNetworkIDFieldName = facebookSocialIDFieldName
	}

	if socialNetworkType == models.SocialNetworkTypeGoogle {
		socialNetworkIDFieldName = googleSocialIDFieldName
	}

	queryBuilder := r.psql.Select(fullIdentityFields...).
		From(identitiesTableName).
		Where(sq.Eq{socialNetworkIDFieldName: socialID})

	row, err := r.queryRow(executor, queryBuilder)

	if err != nil {
		return nil, err
	}

	i, err := r.scanFullIdentity(row)

	if err != nil {
		switch errors.Cause(err) {
		case sql.ErrNoRows:
			return nil, database.ErrRepositoryNoEntitiesFound
		default:
			return nil, errors.Wrap(err, repositoryPrefix+"GetIdentityBySocialID error")
		}
	}

	return i, nil
}

func (r *repository) GetIdentitiesByAccountID(
	executor database.QueryExecutor,
	accountID uuid.UUID,
) ([]*models.Identity, error) {
	queryBuilder := r.psql.Select(fullIdentityFields...).
		From(identitiesTableName).
		Where(sq.Eq{accountIDFieldName: accountID})

	identities := make([]*models.Identity, 0)

	err := r.query(executor, queryBuilder, func(rows *sql.Rows) error {
		identity, scanErr := r.scanFullIdentity(rows)

		if scanErr != nil {
			return errors.Wrap(scanErr, repositoryPrefix+"can`t scan info for GetIdentitiesByAccountID")
		}

		identities = append(identities, identity)
		return nil
	})

	if err != nil {
		return nil, err
	}

	if len(identities) == 0 {
		return nil, database.ErrRepositoryNoEntitiesFound
	}

	return identities, nil
}
