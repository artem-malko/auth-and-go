package service

import (
	"github.com/artem-malko/auth-and-go/services/account"
)

type accountService struct {
	accountRepository account.Repository
}

func New(accountRepository account.Repository) account.Service {
	return &accountService{
		accountRepository: accountRepository,
	}
}
