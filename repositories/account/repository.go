package account

import (
	"errors"

	"github.com/artem-malko/auth-and-go/infrastructure/database"

	"github.com/artem-malko/auth-and-go/models"
	"github.com/google/uuid"
)

var (
	ErrRepositoryAccountNameConstraint = errors.New("account with the same name already exists")
)

// Repository is an interface for any repository
type Repository interface {
	GetAccountByID(executor database.QueryExecutor, accountID uuid.UUID) (*models.Account, error)
	GetAccountsByIDsList(executor database.QueryExecutor, ids []uuid.UUID) ([]*models.Account, error)
	GetAccountByName(executor database.QueryExecutor, accountName string) (*models.Account, error)
	CreateAccount(executor database.QueryExecutor, account models.Account) (*models.Account, error)
	LoginAccount(executor database.QueryExecutor, account uuid.UUID, clientIP string) (*models.Account, error)
	ConfirmAccount(executor database.QueryExecutor, accountID uuid.UUID) (*models.Account, error)
	UpdateAccountName(executor database.QueryExecutor, accountID uuid.UUID, name string) error
	DeactivateAccountByID(executor database.QueryExecutor, accountID uuid.UUID) error
	DeleteUnconfirmedAccountsByAccountIDs(executor database.QueryExecutor, accountIDs []uuid.UUID) error
}
