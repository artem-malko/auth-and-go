package service

import (
	"github.com/artem-malko/auth-and-go/infrastructure/database"
	"github.com/artem-malko/auth-and-go/models"
	"github.com/artem-malko/auth-and-go/services/account"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

func (s *accountService) GetAccountByID(
	executor database.QueryExecutor,
	accountID uuid.UUID,
) (*models.Account, error) {
	accountByID, err := s.accountRepository.GetAccountByID(executor, accountID)

	if err != nil {
		switch errors.Cause(err) {
		case database.ErrRepositoryNoEntitiesFound:
			return nil, account.ErrAccountNotFound
		default:
			return nil, errors.Wrap(err, "account service: GetAccountByID error")
		}
	}

	return accountByID, nil
}

func (s *accountService) GetAccountByName(
	executor database.QueryExecutor,
	accountName string,
) (*models.Account, error) {
	accountByName, err := s.accountRepository.GetAccountByName(executor, accountName)

	if err != nil {
		switch errors.Cause(err) {
		case database.ErrRepositoryNoEntitiesFound:
			return nil, account.ErrAccountNotFound
		default:
			return nil, errors.Wrap(err, "account service: GetAccountByName error")
		}
	}

	return accountByName, nil
}
