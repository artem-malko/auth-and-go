package service

import (
	"github.com/artem-malko/auth-and-go/infrastructure/database"
	"github.com/artem-malko/auth-and-go/services/account"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

func (s *accountService) DeactivateAccountByID(
	executor database.QueryExecutor,
	accountID uuid.UUID,
) (err error) {
	err = s.accountRepository.DeactivateAccountByID(executor, accountID)

	if err != nil {
		switch errors.Cause(err) {
		case database.ErrRepositoryNoRowsAffected:
			return account.ErrNoAccountsUpdated
		}
		return errors.Wrap(err, "account service: DeactivateAccountByID")
	}

	return nil
}

func (s *accountService) DeleteUnconfirmedAccountsByAccountIDs(
	executor database.QueryExecutor,
	accountIDs []uuid.UUID,
) error {
	err := s.accountRepository.DeleteUnconfirmedAccountsByAccountIDs(executor, accountIDs)

	return errors.Wrap(err, "account service: DeleteUnconfirmedAccountsByAccountIDs")
}
