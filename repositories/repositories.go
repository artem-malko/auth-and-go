package repositories

import (
	"database/sql"

	"github.com/pkg/errors"

	"github.com/artem-malko/auth-and-go/infrastructure/database"
	migrate "github.com/rubenv/sql-migrate"

	"github.com/apex/log"

	"github.com/artem-malko/auth-and-go/repositories/account"
	accountRepository "github.com/artem-malko/auth-and-go/repositories/account/repository"
	"github.com/artem-malko/auth-and-go/repositories/identity"
	identityRepository "github.com/artem-malko/auth-and-go/repositories/identity/repository"
	"github.com/artem-malko/auth-and-go/repositories/session"
	sessionRepository "github.com/artem-malko/auth-and-go/repositories/session/repository"
	"github.com/artem-malko/auth-and-go/repositories/token"
	tokenRepository "github.com/artem-malko/auth-and-go/repositories/token/repository"
)

type Repositories struct {
	SessionRepository  session.Repository
	AccountRepository  account.Repository
	IdentityRepository identity.Repository
	TokenRepository    token.Repository
}

func New(db *sql.DB, logger log.Interface) (*Repositories, error) {
	err := RunMigrations(db, logger)

	if err != nil {
		return nil, errors.Wrap(err, "repositories: run migrations error")
	}

	accountRepositoryInstance, err := accountRepository.New()

	if err != nil {
		return nil, errors.Wrap(err, "repositories: can`t create account repository")
	}

	identityRepositoryInstance, err := identityRepository.New()

	if err != nil {
		return nil, errors.Wrap(err, "repositories: can`t create identity repository")
	}

	sessionRepositoryInstance, err := sessionRepository.New()

	if err != nil {
		return nil, errors.Wrap(err, "repositories: can`t create session repository")
	}

	tokenRepositoryInstance, err := tokenRepository.New()

	if err != nil {
		return nil, errors.Wrap(err, "repositories: can`t create token repository")
	}

	return &Repositories{
		AccountRepository:  accountRepositoryInstance,
		IdentityRepository: identityRepositoryInstance,
		SessionRepository:  sessionRepositoryInstance,
		TokenRepository:    tokenRepositoryInstance,
	}, nil
}

func RunMigrations(db *sql.DB, logger log.Interface) error {
	migrate.SetTable(database.MigrationsTableName)

	migrations := &migrate.MemoryMigrationSource{
		Migrations: make([]*migrate.Migration, 0),
	}

	migrations.Migrations = append(migrations.Migrations, accountRepository.Migrations...)
	migrations.Migrations = append(migrations.Migrations, sessionRepository.Migrations...)
	migrations.Migrations = append(migrations.Migrations, identityRepository.Migrations...)
	migrations.Migrations = append(migrations.Migrations, tokenRepository.Migrations...)

	migrationsCount, err := migrate.Exec(db, "postgres", migrations, migrate.Up)

	logger.Infof("%d migrations are applied", migrationsCount)

	return err
}
