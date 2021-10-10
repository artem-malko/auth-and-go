package account

import (
	"errors"

	"github.com/artem-malko/auth-and-go/infrastructure/database"
	"github.com/artem-malko/auth-and-go/models"
	"github.com/google/uuid"
)

var (
	ErrAccountNotFound   = errors.New("account not found")
	ErrNoAccountsUpdated = errors.New("accounts with sent params are not updated")
	ErrAccountNameExists = errors.New("account with sent name exists")
)

type Repository interface {
	GetAccountByID(executor database.QueryExecutor, accountID uuid.UUID) (*models.Account, error)
	GetAccountByName(executor database.QueryExecutor, accountName string) (*models.Account, error)
	CreateAccount(executor database.QueryExecutor, account models.Account) (*models.Account, error)
	LoginAccount(executor database.QueryExecutor, account uuid.UUID, clientIP string) (*models.Account, error)
	ConfirmAccount(executor database.QueryExecutor, accountID uuid.UUID) (*models.Account, error)
	UpdateAccountName(executor database.QueryExecutor, accountID uuid.UUID, name string) error
	DeactivateAccountByID(executor database.QueryExecutor, accountID uuid.UUID) error
	DeleteUnconfirmedAccountsByAccountIDs(
		executor database.QueryExecutor,
		accountIDs []uuid.UUID,
	) error
}

type Service interface {
	GetAccountByID(executor database.QueryExecutor, accountID uuid.UUID) (*models.Account, error)
	GetAccountByName(executor database.QueryExecutor, accountName string) (*models.Account, error)
	LoginAccount(
		executor database.QueryExecutor,
		account uuid.UUID,
		clientIP string,
	) (*models.Account, error)
	CreateAccount(
		executor database.QueryExecutor,
		email string,
		accountStatus models.AccountStatus,
	) (*models.Account, error)
	CreateAccountWithOAuth(
		executor database.QueryExecutor,
		email, firstName, lastName, avatarURL string,
	) (*models.Account, error)
	DeactivateAccountByID(executor database.QueryExecutor, accountID uuid.UUID) error
	DeleteUnconfirmedAccountsByAccountIDs(executor database.QueryExecutor, accountIDs []uuid.UUID) error
	UpdateAccountName(executor database.QueryExecutor, accountID uuid.UUID, name string) error
	ConfirmAccount(executor database.QueryExecutor, accountID uuid.UUID) (*models.Account, error)
}
