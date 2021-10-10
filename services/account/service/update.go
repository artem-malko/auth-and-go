package service

import (
	"github.com/artem-malko/auth-and-go/infrastructure/database"
	"github.com/artem-malko/auth-and-go/models"
	accountRepository "github.com/artem-malko/auth-and-go/repositories/account"
	"github.com/artem-malko/auth-and-go/services/account"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

func (s *accountService) UpdateAccountName(
	executor database.QueryExecutor,
	accountID uuid.UUID,
	name string,
) (err error) {
	err = s.accountRepository.UpdateAccountName(executor, accountID, name)

	if err != nil {
		switch errors.Cause(err) {
		case database.ErrRepositoryNoRowsAffected:
			return account.ErrNoAccountsUpdated
		case accountRepository.ErrRepositoryAccountNameConstraint:
			return account.ErrAccountNameExists
		default:
			return errors.Wrap(err, "account service: UpdateAccountName error")
		}
	}

	return nil
}

func (s *accountService) ConfirmAccount(
	executor database.QueryExecutor,
	accountID uuid.UUID,
) (*models.Account, error) {
	updatedAccount, err := s.accountRepository.ConfirmAccount(executor, accountID)

	if err != nil {
		switch errors.Cause(err) {
		case database.ErrRepositoryNoRowsAffected:
			return nil, account.ErrNoAccountsUpdated
		default:
			return nil, errors.Wrap(err, "account service: ConfirmAccount error")
		}
	}

	return updatedAccount, err
}
