package service

import (
	"github.com/artem-malko/auth-and-go/infrastructure/database"
	"github.com/artem-malko/auth-and-go/services/account"

	"github.com/pkg/errors"

	"github.com/artem-malko/auth-and-go/models"
	"github.com/google/uuid"
)

func (s *accountService) LoginAccount(
	executor database.QueryExecutor,
	accountID uuid.UUID,
	clientIP string,
) (*models.Account, error) {
	a, err := s.accountRepository.LoginAccount(executor, accountID, clientIP)

	if err != nil {
		switch errors.Cause(err) {
		case database.ErrRepositoryNoRowsAffected:
			return nil, account.ErrNoAccountsUpdated
		default:
			return nil, errors.Wrap(err, "account service: login account error")
		}
	}

	return a, nil
}
